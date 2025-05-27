import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { dbCookie } from "./cookies";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function isClient() {
  return typeof window !== "undefined";
}

export function isServer() {
  return !isClient();
}

export function getDbIdFromCookies(cookies: string): Promise<string> {
  return dbCookie.parse(cookies);
}
