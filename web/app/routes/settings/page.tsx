import { AppHeader } from "~/components/shared/header";
import type { Route } from "./+types/page";
import { getServerAuthToken } from "~/lib/auth";
import { initializeServices } from "~/services";
import { redirect, useFetcher, useRevalidator } from "react-router";
import { organizationCookie } from "~/lib/cookies";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { RefreshButton } from "~/components/shared/refresh-button";
import { Save, RefreshCw, Database, Settings, Mail, Cloud, RotateCcw } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { Switch } from "~/components/ui/switch";
import { Textarea } from "~/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";
import type { SettingsData } from "~/services/settings";
import {MailSettingsCard} from "~/routes/settings/mail-card";
import {StorageSettingsCard} from "~/routes/settings/storage-card";
import {ApplicationSettingsCard} from "~/routes/settings/application-card";
import {ApiSettingsCard} from "~/routes/settings/api-card";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Settings - Fluxend" },
    { name: "description", content: "Manage your account settings" },
  ];
}

export async function loader({ request }: Route.LoaderArgs) {
  const authToken = await getServerAuthToken(request.headers);

  if (!authToken) {
    return redirect("/logout");
  }

  const organizationId = await organizationCookie.parse(
      request.headers.get("Cookie")
  );

  if (!organizationId) {
    return redirect("/logout");
  }

  const services = initializeServices(authToken);

  try {
    const settings = await services.settings.getSettings();
    return { settings, error: null };
  } catch (error) {
    console.error("Failed to load settings:", error);
    return {
      settings: null,
      error: error instanceof Error ? error.message : "Failed to load settings"
    };
  }
}

export async function action({ request }: Route.ActionArgs) {
  const authToken = await getServerAuthToken(request.headers);

  if (!authToken) {
    return redirect("/logout");
  }

  const formData = await request.formData();
  const actionType = formData.get("_action") as string;

  const services = initializeServices(authToken);

  try {
    if (actionType === "reset") {
      // Reset settings
      await services.settings.resetSettings();
      return { ok: true, message: "Settings reset successfully", reset: true };
    } else {
      // Update settings
      const settings = Object.fromEntries(formData);
      delete settings._action; // Remove the action type from settings

      // Process form data with proper type conversion
      const processedSettings: SettingsData = {
        appTitle: settings.appTitle as string,
        appUrl: settings.appUrl as string,
        jwtSecret: settings.jwtSecret as string,
        storageDriver: settings.storageDriver as string,
        mailDriver: settings.mailDriver as string,
        maxProjectsPerOrg: parseInt(settings.maxProjectsPerOrg as string) || 10,
        allowRegistrations: settings.allowRegistrations === "on",
        allowProjects: settings.allowProjects === "on",
        allowForms: settings.allowForms === "on",
        allowStorage: settings.allowStorage === "on",
        allowBackups: settings.allowBackups === "on",
        storageMaxContainers: parseInt(settings.storageMaxContainers as string) || 100,
        storageMaxFileSizeInKB: parseInt(settings.storageMaxFileSizeInKB as string) || 10240,
        storageAllowedMimes: settings.storageAllowedMimes as string,
        apiThrottleLimit: parseInt(settings.apiThrottleLimit as string) || 100,
        apiThrottleInterval: parseInt(settings.apiThrottleInterval as string) || 60,
        allowApiThrottle: settings.allowApiThrottle === "on",
        awsAccessKeyId: settings.awsAccessKeyId as string,
        awsSecretAccessKey: settings.awsSecretAccessKey as string,
        awsRegion: settings.awsRegion as string,
        backblazeKeyId: settings.backblazeKeyId as string,
        backblazeApplicationKey: settings.backblazeApplicationKey as string,
        dropboxAccessToken: settings.dropboxAccessToken as string,
        dropboxAppKey: settings.dropboxAppKey as string,
        sendgridApiKey: settings.sendgridApiKey as string,
        sendgridEmailSource: settings.sendgridEmailSource as string,
        mailgunApiKey: settings.mailgunApiKey as string,
        mailgunEmailSource: settings.mailgunEmailSource as string,
        mailgunDomain: settings.mailgunDomain as string,
        mailgunRegion: settings.mailgunRegion as string,
      };

      await services.settings.updateSettings(processedSettings);
      return { ok: true, message: "Settings updated successfully" };
    }
  } catch (error) {
    console.error("Settings action error:", error);
    return {
      ok: false,
      error: error instanceof Error ? error.message : "Failed to process settings"
    };
  }
}

const SettingsPage = ({ loaderData }: Route.ComponentProps) => {
  const { settings: initialSettings, error: loadError } = loaderData;
  const revalidator = useRevalidator();
  const fetcher = useFetcher();
  const [formData, setFormData] = useState<SettingsData | null>(initialSettings);
  const [isDirty, setIsDirty] = useState(false);

  const isSaving = fetcher.state === "submitting";
  const isResetting = fetcher.formData?.get("_action") === "reset";

  useEffect(() => {
    if (fetcher.data?.ok) {
      toast.success(fetcher.data.message);
      setIsDirty(false);

      // If it was a reset action, revalidate to get fresh data
      if (fetcher.data.reset) {
        revalidator.revalidate();
      }
    } else if (fetcher.data?.error) {
      toast.error(fetcher.data.error);
    }
  }, [fetcher.data, revalidator]);

  // Update form data when loader data changes (after reset)
  useEffect(() => {
    if (initialSettings && !isDirty) {
      setFormData(initialSettings);
    }
  }, [initialSettings, isDirty]);

  const handleInputChange = (field: keyof SettingsData, value: any) => {
    if (!formData) return;

    setFormData(prev => prev ? ({
      ...prev,
      [field]: value
    }) : prev);
    setIsDirty(true);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData) return;

    const form = new FormData();
    form.append("_action", "update");

    Object.entries(formData).forEach(([key, value]) => {
      form.append(key, value.toString());
    });

    fetcher.submit(form, { method: "post" });
  };

  const handleReset = () => {
    if (confirm("Are you sure you want to reset all settings to their default values? This action cannot be undone.")) {
      const form = new FormData();
      form.append("_action", "reset");
      fetcher.submit(form, { method: "post" });
    }
  };

  // Show error state if settings failed to load
  if (loadError || !formData) {
    return (
        <div className="min-h-screen bg-gray-50">
          <AppHeader title="Settings">
            <RefreshButton onRefresh={() => revalidator.revalidate()} />
          </AppHeader>
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <Card>
              <CardContent className="py-8">
                <div className="text-center">
                  <p className="text-red-600 mb-4">
                    {loadError || "Failed to load settings"}
                  </p>
                  <Button onClick={() => revalidator.revalidate()}>
                    <RefreshCw className="mr-2 h-4 w-4" />
                    Retry
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
    );
  }

  return (
      <div className="min-h-screen bg-gray-50">
        <AppHeader title="Settings">
          <div className="flex items-center gap-2">
            <Button
                variant="outline"
                onClick={handleReset}
                disabled={isSaving}
                className="text-red-600 hover:text-red-700"
            >
              {isResetting ? (
                  <>
                    <RefreshCw className="mr-2 h-4 w-4 animate-spin" />
                    Resetting...
                  </>
              ) : (
                  <>
                    <RotateCcw className="mr-2 h-4 w-4" />
                    Reset to Defaults
                  </>
              )}
            </Button>
            <RefreshButton
                onRefresh={() => {
                  revalidator.revalidate();
                }}
            />
          </div>
        </AppHeader>

        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <fetcher.Form method="post" onSubmit={handleSubmit} className="space-y-8">
            {/* Application Settings & API Settings */}
            <div className="grid grid-cols-1 xl:grid-cols-2 gap-8">
              {/* Application Settings */}
              <ApplicationSettingsCard formData={formData} onInputChange={handleInputChange}></ApplicationSettingsCard>

              {/* Storage Settings */}
              <StorageSettingsCard formData={formData} onInputChange={handleInputChange} ></StorageSettingsCard>
            </div>

            {/* Storage Settings & Mail Settings */}
            <div className="grid grid-cols-1 xl:grid-cols-2 gap-8">
              {/* API Settings */}
              <ApiSettingsCard formData={formData} onInputChange={handleInputChange}></ApiSettingsCard>

              {/* Mail Settings */}
              <MailSettingsCard formData={formData} onInputChange={handleInputChange}></MailSettingsCard>
            </div>

            {/* Save Button */}
            <div className="flex justify-end">
              <Button
                  type="submit"
                  disabled={!isDirty || isSaving}
                  className="px-8 py-2 cursor-pointer"
                  size="lg"
              >
                {isSaving && !isResetting ? (
                    <>
                      <RefreshCw className="mr-2 h-4 w-4 animate-spin" />
                      Saving Settings...
                    </>
                ) : (
                    <>
                      <Save className="mr-2 h-4 w-4" />
                      Save Settings
                    </>
                )}
              </Button>
            </div>
          </fetcher.Form>
        </div>
      </div>
  );
};

export default SettingsPage;