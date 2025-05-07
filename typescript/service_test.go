package typescript_test

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/lunagic/typescript-go/typescript"
)

type CustomTime time.Time

func (c CustomTime) String() string {
	return time.Time(c).String()
}

func (c CustomTime) MarshalJSON() ([]byte, error) {
	return time.Time(c).MarshalJSON()
}

type TestUser struct {
	Username string
}

func TestPrimary(t *testing.T) {
	type UserID uint64

	type Group struct {
		Name      string `json:"groupName"`
		UpdatedAt time.Time
		DeletedAt *time.Time
		Timeout   time.Duration
		CreateAt  CustomTime
		Data      any
		MoreData  interface{}
	}

	type BaseType struct {
		ID uint64
	}

	type ExtendedType struct {
		BaseType
		Name string
	}

	type TypeNotGivenToTheRegistry struct{}

	type StringTypeNotGivenToTheRegistry string

	type GroupMap map[string]Group

	type User struct {
		Reports           map[UserID]bool
		UserID            UserID `json:"userID"`
		PrimaryGroup      Group  `json:"primaryGroup"`
		UnknownType       TypeNotGivenToTheRegistry
		UnknownStringType StringTypeNotGivenToTheRegistry
		SecondaryGroup    *Group   `json:"secondaryGroup,omitempty"`
		Tags              []string `json:"user_tags"`
		Private           any      `json:"-"`
		unexported        any
	}

	type BaseResponse[T any] struct {
		UpdatedAt time.Time `json:"updated_at"`
		GroupMap  GroupMap  `json:"group_map"`
		Data      T         `json:"data"`
		DataPtr   *T        `json:"data_ptr"`
	}

	_ = User{}.unexported

	type UsersResponse BaseResponse[[]User]

	service := typescript.New(
		typescript.WithTypes(map[string]reflect.Type{
			"TestUserID":    reflect.TypeFor[UserID](),
			"GroupResponse": reflect.TypeFor[BaseResponse[Group]](),
			"UserResponse":  reflect.TypeFor[UsersResponse](),
			"group":         reflect.TypeFor[Group](),
			"SystemUser":    reflect.TypeFor[User](),
			"GroupMapA":     reflect.TypeFor[GroupMap](),
			"GroupMapB":     reflect.TypeFor[map[string]Group](),
			"ExtendedType":  reflect.TypeFor[ExtendedType](),
		}),
		typescript.WithData(map[string]any{
			"foobar": Group{
				Name:      "hello there",
				CreateAt:  CustomTime(time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC)),
				UpdatedAt: time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
			},
		}),
		typescript.WithRoutes(map[string]typescript.Route{
			"userGet": {
				ResponseBody: reflect.TypeFor[UsersResponse](),
				Method:       http.MethodGet,
				Path:         "/api/user",
				QueryParameters: map[string]reflect.Type{
					"userID": reflect.TypeFor[UserID](),
				},
			},
			"userCreate": {
				ResponseBody: reflect.TypeFor[UsersResponse](),
				RequestBody:  reflect.TypeFor[User](),
				Method:       http.MethodPost,
				Path:         "/api/user/create",
			},
		}),
	)

	testThePackage(t, service)
}

func TestSpecialCharacter(t *testing.T) {
	type TestStruct struct {
		Timestamp string `json:"@timestamp"`
		UpdatedAt time.Time
		DeletedAt *time.Time
		Timeout   time.Duration
		Data      any
		MoreData  interface{}
	}

	service := typescript.New(
		typescript.WithTypes(map[string]reflect.Type{
			"TestStruct": reflect.TypeFor[TestStruct](),
		}),
	)

	testThePackage(t, service)
}

func testThePackage(t *testing.T, service *typescript.Service) {
	actualFileName := "test_files/" + t.Name() + "_actual.ts"

	actualFile, err := os.Create(actualFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = actualFile.Close()
	}()

	actualFileBuffer := bytes.NewBuffer([]byte{})

	writer := io.MultiWriter(actualFile, actualFileBuffer)

	if err := service.Generate(writer); err != nil {
		t.Fatal(err)
	}

	expectedContents, err := os.ReadFile("test_files/" + t.Name() + "_expected.ts")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actualFileBuffer.Bytes(), expectedContents) {
		wd, _ := os.Getwd()
		t.Fatal("contents don't match: " + wd + "/" + actualFileName)
	}
}
