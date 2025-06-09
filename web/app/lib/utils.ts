import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { dbCookie, organizationCookie } from "./cookies";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function isClient() {
  return typeof window !== "undefined";
}

export function isServer() {
  return !isClient();
}

export function getDbIdCookie(headers: Headers): Promise<string> {
  const cookies = headers.get("cookie");
  return dbCookie.parse(cookies);
}

export function getOrganizationIdCookie(headers: Headers): Promise<string> {
  const cookies = headers.get("cookie");
  return organizationCookie.parse(cookies);
}

// TODO: Validate Response with Types using zod
export async function getTypedResponseData<T>(response: Response): Promise<T> {
  const json = await response.json();
  return { ok: response.ok, status: response.status, ...json };
}

export function formatTimestamp(timestamp: string): {
  date: string;
  time: string;
  fullDate: string;
  relativeTime: string;
} {
  try {
    const date = new Date(timestamp);

    // Check if date is valid
    if (isNaN(date.getTime())) {
      return { date: "Invalid date", time: "", fullDate: "", relativeTime: "" };
    }

    // Format date as "Mon DD, YYYY"
    const monthNames = [
      "Jan",
      "Feb",
      "Mar",
      "Apr",
      "May",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Oct",
      "Nov",
      "Dec",
    ];
    const month = monthNames[date.getMonth()];
    const day = date.getDate();
    const year = date.getFullYear();

    // Format time as "HH:MM AM/PM"
    let hours = date.getHours();
    const minutes = date.getMinutes().toString().padStart(2, "0");
    const seconds = date.getSeconds().toString().padStart(2, "0");
    const ampm = hours >= 12 ? "PM" : "AM";
    hours = hours % 12;
    hours = hours ? hours : 12; // Convert 0 to 12

    // Calculate relative time
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffSec = Math.floor(diffMs / 1000);
    const diffMin = Math.floor(diffSec / 60);
    const diffHour = Math.floor(diffMin / 60);
    const diffDay = Math.floor(diffHour / 24);

    let relativeTime = "";
    if (diffMs < 0) {
      // Future date
      relativeTime = "In the future";
    } else if (diffDay > 365) {
      relativeTime = `${Math.floor(diffDay / 365)} year${
        Math.floor(diffDay / 365) > 1 ? "s" : ""
      } ago`;
    } else if (diffDay > 30) {
      relativeTime = `${Math.floor(diffDay / 30)} month${
        Math.floor(diffDay / 30) > 1 ? "s" : ""
      } ago`;
    } else if (diffDay > 0) {
      relativeTime = `${diffDay} day${diffDay > 1 ? "s" : ""} ago`;
    } else if (diffHour > 0) {
      relativeTime = `${diffHour} hour${diffHour > 1 ? "s" : ""} ago`;
    } else if (diffMin > 0) {
      relativeTime = `${diffMin} minute${diffMin > 1 ? "s" : ""} ago`;
    } else {
      relativeTime = "Just now";
    }

    // Full ISO date for tooltip
    const fullDate = `${year}-${(date.getMonth() + 1)
      .toString()
      .padStart(2, "0")}-${day.toString().padStart(2, "0")} ${hours
      .toString()
      .padStart(2, "0")}:${minutes}:${seconds} ${ampm}`;

    return {
      date: `${month} ${day}, ${year}`,
      time: `${hours}:${minutes} ${ampm}`,
      fullDate,
      relativeTime,
    };
  } catch (error) {
    return { date: "Invalid date", time: "", fullDate: "", relativeTime: "" };
  }
}
