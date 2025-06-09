import { stringify } from "qs";
import { isServer } from "../lib/utils";

export type APIRequestOptions = {
  baseUrl?: string;
  params?: Record<string, any>;
  timeout?: number;
  [key: string]: any; // Allow additional fetch options
} & RequestInit;

// Get the base URL from the appropriate environment variable
const getBaseUrl = (): string => {
  if (isServer()) {
    const serverUrl = process.env.VITE_FLX_INTERNAL_URL;
    if (serverUrl) {
      return serverUrl;
    }

    console.warn("VITE_FLX_INTERNAL_URL environment variable not set");
    return "";
  }

  // Client-side environment variables - check if process exists first
  let clientUrl: string | undefined;

  if (typeof process !== "undefined" && process.env?.VITE_FLX_API_URL) {
    clientUrl = process.env.VITE_FLX_API_URL;
  } else if (
    typeof import.meta !== "undefined" &&
    import.meta.env?.VITE_FLX_API_URL
  ) {
    clientUrl = import.meta.env.VITE_FLX_API_URL;
  }

  if (clientUrl) {
    return clientUrl;
  }

  console.warn("No VITE_FLX_API_URL found in environment variables");

  return "";
};

/**
 * Client-side fetch implementation using browser's fetch API
 */
const clientFetch = async (
  url: string,
  method: string,
  data?: any,
  options: APIRequestOptions = {}
): Promise<Response> => {
  const { params, headers, baseUrl: customBaseUrl, ...restOptions } = options;
  const baseUrl = customBaseUrl ? customBaseUrl : getBaseUrl();

  let fullUrl = `${baseUrl}${url}`;

  // Add query parameters if they exist
  if (params) {
    const queryString = stringify(params);
    fullUrl = `${fullUrl}${fullUrl.includes("?") ? "&" : "?"}${queryString}`;
  }

  const fetchOptions: RequestInit = {
    method,
    headers: {
      "Content-Type": "application/json",
      ...headers,
    },
    ...restOptions,
  };

  // Add body for methods that support it
  if (data && ["POST", "PUT", "PATCH"].includes(method)) {
    fetchOptions.body = JSON.stringify(data);
  }

  return fetch(fullUrl, fetchOptions);
};

/**
 * Server-side fetch implementation using Node.js native fetch
 */
const serverFetch = async (
  url: string,
  method: string,
  data?: any,
  options: APIRequestOptions = {}
): Promise<Response> => {
  const {
    params,
    headers,
    baseUrl: customBaseUrl,
    timeout = 30000,
    ...restOptions
  } = options;
  const baseUrl = customBaseUrl ? customBaseUrl : getBaseUrl();

  let fullUrl = `${baseUrl}${url}`;

  // Add query parameters if they exist
  if (params) {
    const queryString = stringify(params);
    fullUrl = `${fullUrl}${fullUrl.includes("?") ? "&" : "?"}${queryString}`;
  }

  const controller = new AbortController();

  // Set up timeout
  setTimeout(() => {
    controller.abort();
  }, timeout);

  const fetchOptions: RequestInit = {
    method,
    headers: {
      "Content-Type": "application/json",
      ...headers,
    },
    signal: controller.signal,
    ...restOptions,
  };

  // Add body for methods that support it
  if (data && ["POST", "PUT", "PATCH"].includes(method)) {
    fetchOptions.body = JSON.stringify(data);
  }

  return fetch(fullUrl, fetchOptions);
};

/**
 * Universal fetch function that uses the appropriate implementation based on environment
 */
const universalFetch = (
  url: string,
  method: string,
  data?: any,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return isServer()
    ? serverFetch(url, method, data, options)
    : clientFetch(url, method, data, options);
};

/**
 * HTTP GET request
 * @param url - API endpoint to call
 * @param options - Request options including optional query params
 * @returns Promise with the parsed JSON response
 */
export const get = (
  url: string,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return universalFetch(url, "GET", undefined, options);
};

/**
 * HTTP POST request
 * @param url - API endpoint to call
 * @param data - Request payload
 * @param options - Request options
 * @returns Promise with the parsed JSON response
 */
export const post = (
  url: string,
  data: any,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return universalFetch(url, "POST", data, options);
};

/**
 * HTTP PUT request
 * @param url - API endpoint to call
 * @param data - Request payload
 * @param options - Request options
 * @returns Promise with the parsed JSON response
 */
export const put = (
  url: string,
  data: any,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return universalFetch(url, "PUT", data, options);
};

/**
 * HTTP PATCH request
 * @param url - API endpoint to call
 * @param data - Request payload
 * @param options - Request options
 * @returns Promise with the parsed JSON response
 */
export const patch = (
  url: string,
  data: any,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return universalFetch(url, "PATCH", data, options);
};

/**
 * HTTP DELETE request
 * @param url - API endpoint to call
 * @param options - Request options
 * @returns Promise with the parsed JSON response
 */
export const del = (
  url: string,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return universalFetch(url, "DELETE", undefined, options);
};

export default {
  get,
  post,
  put,
  patch,
  del,
};
