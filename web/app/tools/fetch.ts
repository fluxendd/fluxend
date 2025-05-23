import { stringify } from "qs";
import { isServer } from "../lib/utils";

export type APIRequestOptions = {
  baseUrl?: string;
  params?: Record<string, any>;
  timeout?: number;
  [key: string]: any; // Allow additional fetch options
} & RequestInit;

// Get the base URL from the appropriate environment variable
const getBaseUrl = () => {
  if (isServer()) {
    return process.env.VITE_FLX_API_BASE_URL || "";
  } else {
    return import.meta.env.VITE_FLX_API_BASE_URL || "";
  }
};

/**
 * Client-side fetch implementation using browser's fetch API
 */
const clientFetch = async <T>(
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
const serverFetch = async <T>(
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
const universalFetch = <T>(
  url: string,
  method: string,
  data?: any,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return isServer()
    ? serverFetch<T>(url, method, data, options)
    : clientFetch<T>(url, method, data, options);
};

/**
 * HTTP GET request
 * @param url - API endpoint to call
 * @param options - Request options including optional query params
 * @returns Promise with the parsed JSON response
 */
export const get = <T>(
  url: string,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return universalFetch<T>(url, "GET", undefined, options);
};

/**
 * HTTP POST request
 * @param url - API endpoint to call
 * @param data - Request payload
 * @param options - Request options
 * @returns Promise with the parsed JSON response
 */
export const post = <T>(
  url: string,
  data: any,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return universalFetch<T>(url, "POST", data, options);
};

/**
 * HTTP PUT request
 * @param url - API endpoint to call
 * @param data - Request payload
 * @param options - Request options
 * @returns Promise with the parsed JSON response
 */
export const put = <T>(
  url: string,
  data: any,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return universalFetch<T>(url, "PUT", data, options);
};

/**
 * HTTP PATCH request
 * @param url - API endpoint to call
 * @param data - Request payload
 * @param options - Request options
 * @returns Promise with the parsed JSON response
 */
export const patch = <T>(
  url: string,
  data: any,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return universalFetch<T>(url, "PATCH", data, options);
};

/**
 * HTTP DELETE request
 * @param url - API endpoint to call
 * @param options - Request options
 * @returns Promise with the parsed JSON response
 */
export const del = <T>(
  url: string,
  options: APIRequestOptions = {}
): Promise<Response> => {
  return universalFetch<T>(url, "DELETE", undefined, options);
};

export default {
  get,
  post,
  put,
  patch,
  del,
};
