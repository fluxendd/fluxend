import { redirect } from "react-router";
import { sessionCookie, organizationCookie, projectCookie } from "../../lib/cookies";

export async function loader() {
  const sessionCookieStr = await sessionCookie.serialize("");
  const organizationCookieStr = await organizationCookie.serialize("");
  const projectCookieStr = await projectCookie.serialize("");

  return redirect("/", {
    headers: {
      "Set-Cookie": [
        sessionCookieStr,
        organizationCookieStr,
        projectCookieStr
      ].join(", "),
    },
  });
}
