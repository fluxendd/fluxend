import { createCookie } from "react-router";

export const sessionCookie = createCookie("session");
export const organizationCookie = createCookie("organization_uuid");
export const projectCookie = createCookie("project_uuid");
export const dbCookie = createCookie("db_uuid");
