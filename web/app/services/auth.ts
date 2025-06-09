import type { APIResponse } from "~/lib/types";
import { getTypedResponseData } from "~/lib/utils";
import { post } from "~/tools/fetch";

// NOTE: Auth Service is special kind of service that does not rely on the auth token.
// That's why it should not be initialized through initilizeService function.

export const login = async (email: string, password: string) => {
  const response = await post("/users/login", {
    email,
    password,
  });

  const data = await getTypedResponseData<
    APIResponse<{ token: string; organizationUuid: string }>
  >(response);
  return data;
};

export const signup = async (
  email: string,
  username: string,
  password: string
) => {
  const response = await post("/users/register", {
    email,
    username,
    password,
  });

  const data = await getTypedResponseData<
    APIResponse<{ token: string; organizationUuid: string }>
  >(response);
  return data;
};
