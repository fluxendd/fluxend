import { sessionCookie, organizationCookie } from "../../lib/cookies";

export async function loader() {
  const sessionCookieStr = await sessionCookie.serialize("");
  const organizationCookieStr = await organizationCookie.serialize("");

  // Create a Response object that allows setting multiple cookies
  const response = new Response(null, {
    status: 302,
    headers: {
      Location: "/",
    },
  });

  // Append each cookie separately
  response.headers.append("Set-Cookie", sessionCookieStr);
  response.headers.append("Set-Cookie", organizationCookieStr);

  return response;
}
