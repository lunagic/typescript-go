// NOTE: This file was auto-generated
// and should NOT be edited manually.

export namespace GoGenerated {
	export type ExtendedType = {
		ID: number
		Name: string
	}

	export type GroupMapA = { [key: string]: group } | null

	export type GroupMapB = { [key: string]: group } | null

	export type GroupResponse = {
		updated_at: string
		group_map: { [key: string]: group } | null
		data: group
		data_ptr: group | null
	}

	export type SystemUser = {
		Reports: { [key: TestUserID]: boolean } | null
		userID: TestUserID
		primaryGroup: group
		UnknownType: unknown
		UnknownStringType: string
		secondaryGroup?: group | null
		user_tags: string[] | null
	}

	export type TestUserID = number

	export type UserResponse = {
		updated_at: string
		group_map: { [key: string]: group } | null
		data: SystemUser[] | null
		data_ptr: SystemUser[] | null
	}

	export type group = {
		groupName: string
		UpdatedAt: string
		DeletedAt: string | null
		Timeout: number
		CreateAt: string
		Data: any
		MoreData: any
	}

	export const userCreate = async (payload: SystemUser) => {
		const response = await fetch("/api/user/create", {
			method: "POST",
			body: JSON.stringify(payload),
		})

		return await response.json() as UserResponse
	}

	export const userGet = async (userID: TestUserID) => {
		const params = {
			userID: userID,
		}

		const queryString = Object.keys(params).map((key) => {
			return encodeURIComponent(key) + "=" + encodeURIComponent(params[key])
		}).join("&")

		const response = await fetch(`/api/user?${queryString}`, {
			method: "GET",
		})

		return await response.json() as UserResponse
	}

	export const foobar: group = {
		"groupName": "hello there",
		"UpdatedAt": "0001-01-01T01:01:01.000000001Z",
		"DeletedAt": null,
		"Timeout": 0,
		"CreateAt": "0001-01-01T01:01:01.000000001Z",
		"Data": null,
		"MoreData": null
	}
}
