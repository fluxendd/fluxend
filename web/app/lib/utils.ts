import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { dbCookie } from "./cookies";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function isServer() {
  return typeof window === "undefined";
}

export function isClient() {
  return typeof window !== "undefined";
}

export function getDbIdFromCookies(cookies: string): Promise<string> {
  return dbCookie.parse(cookies);
}
