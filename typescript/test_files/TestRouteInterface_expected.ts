// NOTE: This file was auto-generated
// and should NOT be edited manually.

export namespace GoGenerated {
	export type TestStruct = {
		"@timestamp": string
		UpdatedAt: string
		DeletedAt: string | null
		Timeout: number
		Data: any
		MoreData: any
	}

	export type TestUser = {
		Username: string
	}

	export const GetThing = async (payload: string) => {
		const response = await fetch("/_backend?method=GetThing", {
			method: "POST",
			body: JSON.stringify(payload),
		})

		return await response.json() as TestUser
	}
}
