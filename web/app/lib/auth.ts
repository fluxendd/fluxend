import { isClient, isServer } from "./utils";
import { sessionCookie } from "./cookies";

/**
 * Gets the authentication token from the session cookie
 * Works in both client and server environments
 */
export const getAuthToken = async (
  headers: Headers
): Promise<string | null> => {
  try {
    if (isServer()) {
      const cookieHeader = headers.get("Cookie");
      const sessionToken = await sessionCookie.parse(cookieHeader);
      return sessionToken;
    } else {
      const cookie = await sessionCookie.parse(document.cookie);
      return cookie;
    }
  } catch (error) {
    console.error("Error retrieving auth token:", error);
  }

  return null;
};
