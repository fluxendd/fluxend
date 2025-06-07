import { post } from "~/tools/fetch";

export type APIResponse<T> = {
  success: boolean;
  errors: string[] | null;
  content: T | null;
};

export const login = async (email: string, password: string) => {
  const response = await post<{ token: string }>("/users/login", {
    email,
    password,
  });
  const data = await response.json();

  return data as APIResponse<any>;
};

export const signup = async (
  email: string,
  username: string,
  password: string
) => {
  const response = await post<{ token: string }>("/users/register", {
    email,
    username,
    password,
  });
  const data = await response.json();
  return data as APIResponse<any>;
};
