export interface Account {
    login:	string,
    password?:	string,
    role?: string
}

export interface AuthRequest {
    login: string,
	password: string,
}

export interface NewUserRequest {
    profile: {
        firstName: string,
        lastName: string,
        email: string,
        login: string,
        mobilePhone: string,
    },
    credentials: {
        password: {
            value: string
        }
    },
}
