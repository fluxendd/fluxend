import { redirect } from "react-router";
import {
  sessionCookie,
  organizationCookie,
  projectCookie,
  dbCookie,
} from "../../lib/cookies";

export async function loader() {
  const sessionCookieStr = await sessionCookie.serialize("");
  const organizationCookieStr = await organizationCookie.serialize("");
  const projectCookieStr = await projectCookie.serialize("");
  const dbCookieStr = await dbCookie.serialize("");

  return redirect("/", {
    headers: {
      "Set-Cookie": [
        sessionCookieStr,
        organizationCookieStr,
        projectCookieStr,
        dbCookieStr,
      ].join(", "),
    },
  });
}
