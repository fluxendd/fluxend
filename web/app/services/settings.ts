import type { APIResponse } from "~/lib/types";
import { getTypedResponseData } from "~/lib/utils";
import { get, put, type APIRequestOptions } from "~/tools/fetch";

export interface SettingResponse {
    id: string;
    name: string;
    value: string;
    defaultValue?: string;
    created_at: string;
    updated_at: string;
}

export interface SettingUpdateRequest {
    settings: Array<{
        name: string;
        value: string;
    }>;
}

export interface SettingsData {
    appTitle: string;
    appUrl: string;
    jwtSecret: string;
    storageDriver: string;
    mailDriver: string;
    maxProjectsPerOrg: number;
    allowRegistrations: boolean;
    allowProjects: boolean;
    allowForms: boolean;
    allowStorage: boolean;
    allowBackups: boolean;
    storageMaxContainers: number;
    storageMaxFileSizeInKB: number;
    storageAllowedMimes: string;
    apiThrottleLimit: number;
    apiThrottleInterval: number;
    allowApiThrottle: boolean;
    awsAccessKeyId: string;
    awsSecretAccessKey: string;
    awsRegion: string;
    backblazeKeyId: string;
    backblazeApplicationKey: string;
    dropboxAccessToken: string;
    dropboxAppKey: string;
    sendgridApiKey: string;
    sendgridEmailSource: string;
    mailgunApiKey: string;
    mailgunEmailSource: string;
    mailgunDomain: string;
    mailgunRegion: string;
}

export const createSettingsService = (authToken: string) => {
    const getSettings = async (): Promise<SettingsData> => {
        try {
            const fetchOptions: APIRequestOptions = {
                headers: {
                    "Content-Type": "application/json",
                    ...(authToken && { Authorization: `Bearer ${authToken}` }),
                },
            };

            const response = await get("/admin/settings", fetchOptions);
            const result = await getTypedResponseData<APIResponse<SettingResponse[]>>(
                response
            );

            if (result.success && result.content) {
                // Transform array of settings into structured object
                const settingsMap = result.content.reduce((acc, setting) => {
                    acc[setting.name] = setting.value;
                    return acc;
                }, {} as Record<string, string>);


                // Convert string values to appropriate types and provide defaults
                const settings: SettingsData = {
                    appTitle: settingsMap.appTitle || "Fluxend",
                    appUrl: settingsMap.appUrl || "https://app.fluxend.com",
                    jwtSecret: settingsMap.jwtSecret || "",
                    storageDriver: settingsMap.storageDriver || "local",
                    mailDriver: settingsMap.mailDriver || "sendgrid",
                    maxProjectsPerOrg: parseInt(settingsMap.maxProjectsPerOrg) || 10,
                    allowRegistrations: settingsMap.allowRegistrations === "true",
                    allowProjects: settingsMap.allowProjects === "true",
                    allowForms: settingsMap.allowForms === "true",
                    allowStorage: settingsMap.allowStorage === "true",
                    allowBackups: settingsMap.allowBackups === "true",
                    storageMaxContainers: parseInt(settingsMap.storageMaxContainers) || 100,
                    storageMaxFileSizeInKB: parseInt(settingsMap.storageMaxFileSizeInKB) || 10240,
                    storageAllowedMimes: settingsMap.storageAllowedMimes || "image/jpeg,image/png,image/gif,application/pdf,text/plain",
                    apiThrottleLimit: parseInt(settingsMap.apiThrottleLimit) || 100,
                    apiThrottleInterval: parseInt(settingsMap.apiThrottleInterval) || 60,
                    allowApiThrottle: settingsMap.allowApiThrottle === "true",
                    awsAccessKeyId: settingsMap.awsAccessKeyId || "",
                    awsSecretAccessKey: settingsMap.awsSecretAccessKey || "",
                    awsRegion: settingsMap.awsRegion || "us-east-1",
                    backblazeKeyId: settingsMap.backblazeKeyId || "",
                    backblazeApplicationKey: settingsMap.backblazeApplicationKey || "",
                    dropboxAccessToken: settingsMap.dropboxAccessToken || "",
                    dropboxAppKey: settingsMap.dropboxAppKey || "",
                    sendgridApiKey: settingsMap.sendgridApiKey || "",
                    sendgridEmailSource: settingsMap.sendgridEmailSource || "",
                    mailgunApiKey: settingsMap.mailgunApiKey || "",
                    mailgunEmailSource: settingsMap.mailgunEmailSource || "",
                    mailgunDomain: settingsMap.mailgunDomain || "",
                    mailgunRegion: settingsMap.mailgunRegion || "us",
                };

                return settings;
            } else {
                throw new Error(result.errors?.join(", ") || "Unknown error");
            }
        } catch (error) {
            // Fallback to mock data for development
            console.warn("Failed to fetch settings, using mock data:", error);
            return {
                appTitle: "Fluxend",
                appUrl: "https://app.fluxend.com",
                jwtSecret: "your-jwt-secret-here",
                storageDriver: "local",
                mailDriver: "sendgrid",
                maxProjectsPerOrg: 10,
                allowRegistrations: true,
                allowProjects: true,
                allowForms: true,
                allowStorage: true,
                allowBackups: true,
                storageMaxContainers: 100,
                storageMaxFileSizeInKB: 10240,
                storageAllowedMimes: "image/jpeg,image/png,image/gif,application/pdf,text/plain",
                apiThrottleLimit: 100,
                apiThrottleInterval: 60,
                allowApiThrottle: true,
                awsAccessKeyId: "",
                awsSecretAccessKey: "",
                awsRegion: "us-east-1",
                backblazeKeyId: "",
                backblazeApplicationKey: "",
                dropboxAccessToken: "",
                dropboxAppKey: "",
                sendgridApiKey: "",
                sendgridEmailSource: "",
                mailgunApiKey: "",
                mailgunEmailSource: "",
                mailgunDomain: "",
                mailgunRegion: "us",
            };
        }
    };

    const updateSettings = async (settings: SettingsData): Promise<SettingResponse[]> => {
        try {
            const fetchOptions: APIRequestOptions = {
                headers: {
                    "Content-Type": "application/json",
                    ...(authToken && { Authorization: `Bearer ${authToken}` }),
                },
            };

            // Transform SettingsData into the expected API format
            const updateRequest: SettingUpdateRequest = {
                settings: Object.entries(settings).map(([name, value]) => ({
                    name,
                    value: value.toString(),
                })),
            };

            const response = await put("/admin/settings", updateRequest, fetchOptions);
            const result = await getTypedResponseData<APIResponse<SettingResponse[]>>(
                response
            );

            if (result.success && result.content) {
                return result.content;
            } else {
                throw new Error(result.errors?.join(", ") || "Failed to update settings");
            }
        } catch (error) {
            console.error("Failed to update settings:", error);
            throw error;
        }
    };

    const resetSettings = async (): Promise<SettingResponse[]> => {
        try {
            const fetchOptions: APIRequestOptions = {
                headers: {
                    "Content-Type": "application/json",
                    ...(authToken && { Authorization: `Bearer ${authToken}` }),
                },
            };

            const response = await put("/admin/settings/reset", {}, fetchOptions);
            const result = await getTypedResponseData<APIResponse<SettingResponse[]>>(
                response
            );

            if (result.success && result.content) {
                return result.content;
            } else {
                throw new Error(result.errors?.join(", ") || "Failed to reset settings");
            }
        } catch (error) {
            console.error("Failed to reset settings:", error);
            throw error;
        }
    };

    return {
        getSettings,
        updateSettings,
        resetSettings,
    };
};

export type SettingsService = ReturnType<typeof createSettingsService>;