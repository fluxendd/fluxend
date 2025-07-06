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
import { Save, RefreshCw, Database, Shield, Settings, Mail, Cloud, RotateCcw } from "lucide-react";
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

  const storageDriverOptions = [
    { value: "local", label: "Local Storage" },
    { value: "aws", label: "AWS S3" },
    { value: "backblaze", label: "Backblaze B2" },
    { value: "dropbox", label: "Dropbox" },
  ];

  const mailDriverOptions = [
    { value: "sendgrid", label: "SendGrid" },
    { value: "mailgun", label: "Mailgun" },
  ];

  const mailgunRegionOptions = [
    { value: "us", label: "United States" },
    { value: "eu", label: "Europe" },
  ];

  const renderStorageConfiguration = () => {
    switch (formData.storageDriver) {
      case "aws":
        return (
            <div className="space-y-6">
              <div className="flex items-center gap-2 mb-4">
                <Cloud className="h-4 w-4 text-blue-500" />
                <h4 className="font-medium text-sm">AWS S3 Configuration</h4>
              </div>
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="space-y-3">
                  <Label htmlFor="awsAccessKeyId" className="text-sm font-medium">
                    Access Key ID
                  </Label>
                  <Input
                      id="awsAccessKeyId"
                      name="awsAccessKeyId"
                      type="password"
                      value={formData.awsAccessKeyId}
                      onChange={(e) => handleInputChange("awsAccessKeyId", e.target.value)}
                      placeholder="Your AWS Access Key ID"
                      className="mt-2"
                  />
                </div>
                <div className="space-y-3">
                  <Label htmlFor="awsSecretAccessKey" className="text-sm font-medium">
                    Secret Access Key
                  </Label>
                  <Input
                      id="awsSecretAccessKey"
                      name="awsSecretAccessKey"
                      type="password"
                      value={formData.awsSecretAccessKey}
                      onChange={(e) => handleInputChange("awsSecretAccessKey", e.target.value)}
                      placeholder="Your AWS Secret Access Key"
                      className="mt-2"
                  />
                </div>
                <div className="space-y-3">
                  <Label htmlFor="awsRegion" className="text-sm font-medium">
                    AWS Region
                  </Label>
                  <Input
                      id="awsRegion"
                      name="awsRegion"
                      value={formData.awsRegion}
                      onChange={(e) => handleInputChange("awsRegion", e.target.value)}
                      placeholder="us-east-1"
                      className="mt-2"
                  />
                </div>
              </div>
            </div>
        );
      case "backblaze":
        return (
            <div className="space-y-6">
              <div className="flex items-center gap-2 mb-4">
                <Cloud className="h-4 w-4 text-orange-500" />
                <h4 className="font-medium text-sm">Backblaze B2 Configuration</h4>
              </div>
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="space-y-3">
                  <Label htmlFor="backblazeKeyId" className="text-sm font-medium">
                    Key ID
                  </Label>
                  <Input
                      id="backblazeKeyId"
                      name="backblazeKeyId"
                      type="password"
                      value={formData.backblazeKeyId}
                      onChange={(e) => handleInputChange("backblazeKeyId", e.target.value)}
                      placeholder="Your Backblaze Key ID"
                      className="mt-2"
                  />
                </div>
                <div className="space-y-3">
                  <Label htmlFor="backblazeApplicationKey" className="text-sm font-medium">
                    Application Key
                  </Label>
                  <Input
                      id="backblazeApplicationKey"
                      name="backblazeApplicationKey"
                      type="password"
                      value={formData.backblazeApplicationKey}
                      onChange={(e) => handleInputChange("backblazeApplicationKey", e.target.value)}
                      placeholder="Your Backblaze Application Key"
                      className="mt-2"
                  />
                </div>
              </div>
            </div>
        );
      case "dropbox":
        return (
            <div className="space-y-6">
              <div className="flex items-center gap-2 mb-4">
                <Cloud className="h-4 w-4 text-blue-600" />
                <h4 className="font-medium text-sm">Dropbox Configuration</h4>
              </div>
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="space-y-3">
                  <Label htmlFor="dropboxAccessToken" className="text-sm font-medium">
                    Access Token
                  </Label>
                  <Input
                      id="dropboxAccessToken"
                      name="dropboxAccessToken"
                      type="password"
                      value={formData.dropboxAccessToken}
                      onChange={(e) => handleInputChange("dropboxAccessToken", e.target.value)}
                      placeholder="Your Dropbox Access Token"
                      className="mt-2"
                  />
                </div>
                <div className="space-y-3">
                  <Label htmlFor="dropboxAppKey" className="text-sm font-medium">
                    App Key
                  </Label>
                  <Input
                      id="dropboxAppKey"
                      name="dropboxAppKey"
                      type="password"
                      value={formData.dropboxAppKey}
                      onChange={(e) => handleInputChange("dropboxAppKey", e.target.value)}
                      placeholder="Your Dropbox App Key"
                      className="mt-2"
                  />
                </div>
              </div>
            </div>
        );
      default:
        return null;
    }
  };

  const renderMailConfiguration = () => {
    switch (formData.mailDriver) {
      case "sendgrid":
        return (
            <div className="space-y-6">
              <div className="flex items-center gap-2 mb-4">
                <Mail className="h-4 w-4 text-blue-500" />
                <h4 className="font-medium text-sm">SendGrid Configuration</h4>
              </div>
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="space-y-3">
                  <Label htmlFor="sendgridApiKey" className="text-sm font-medium">
                    API Key
                  </Label>
                  <Input
                      id="sendgridApiKey"
                      name="sendgridApiKey"
                      type="password"
                      value={formData.sendgridApiKey}
                      onChange={(e) => handleInputChange("sendgridApiKey", e.target.value)}
                      placeholder="Your SendGrid API Key"
                      className="mt-2"
                  />
                </div>
                <div className="space-y-3">
                  <Label htmlFor="sendgridEmailSource" className="text-sm font-medium">
                    Email Source
                  </Label>
                  <Input
                      id="sendgridEmailSource"
                      name="sendgridEmailSource"
                      type="email"
                      value={formData.sendgridEmailSource}
                      onChange={(e) => handleInputChange("sendgridEmailSource", e.target.value)}
                      placeholder="noreply@yourapp.com"
                      className="mt-2"
                  />
                </div>
              </div>
            </div>
        );
      case "mailgun":
        return (
            <div className="space-y-6">
              <div className="flex items-center gap-2 mb-4">
                <Mail className="h-4 w-4 text-red-500" />
                <h4 className="font-medium text-sm">Mailgun Configuration</h4>
              </div>
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="space-y-3">
                  <Label htmlFor="mailgunApiKey" className="text-sm font-medium">
                    API Key
                  </Label>
                  <Input
                      id="mailgunApiKey"
                      name="mailgunApiKey"
                      type="password"
                      value={formData.mailgunApiKey}
                      onChange={(e) => handleInputChange("mailgunApiKey", e.target.value)}
                      placeholder="Your Mailgun API Key"
                      className="mt-2"
                  />
                </div>
                <div className="space-y-3">
                  <Label htmlFor="mailgunEmailSource" className="text-sm font-medium">
                    Email Source
                  </Label>
                  <Input
                      id="mailgunEmailSource"
                      name="mailgunEmailSource"
                      type="email"
                      value={formData.mailgunEmailSource}
                      onChange={(e) => handleInputChange("mailgunEmailSource", e.target.value)}
                      placeholder="noreply@yourapp.com"
                      className="mt-2"
                  />
                </div>
                <div className="space-y-3">
                  <Label htmlFor="mailgunDomain" className="text-sm font-medium">
                    Domain
                  </Label>
                  <Input
                      id="mailgunDomain"
                      name="mailgunDomain"
                      value={formData.mailgunDomain}
                      onChange={(e) => handleInputChange("mailgunDomain", e.target.value)}
                      placeholder="mg.yourapp.com"
                      className="mt-2"
                  />
                </div>
                <div className="space-y-3">
                  <Label htmlFor="mailgunRegion" className="text-sm font-medium">
                    Region
                  </Label>
                  <Select
                      value={formData.mailgunRegion}
                      onValueChange={(value) => handleInputChange("mailgunRegion", value)}
                  >
                    <SelectTrigger className="mt-2">
                      <SelectValue placeholder="Select region" />
                    </SelectTrigger>
                    <SelectContent>
                      {mailgunRegionOptions.map(option => (
                          <SelectItem key={option.value} value={option.value}>
                            {option.label}
                          </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </div>
        );
      default:
        return null;
    }
  };

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
            {/* Application & Security Settings */}
            <div className="grid grid-cols-1 xl:grid-cols-2 gap-8">
              {/* Application Settings */}
              <Card className="h-fit">
                <CardHeader>
                  <div className="flex items-center gap-2">
                    <Settings className="h-5 w-5 text-blue-500" />
                    <CardTitle>Application Settings</CardTitle>
                  </div>
                  <CardDescription>
                    Configure your application's basic settings
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="space-y-3">
                    <Label htmlFor="appTitle" className="text-sm font-medium">
                      App Title
                    </Label>
                    <Input
                        id="appTitle"
                        name="appTitle"
                        value={formData.appTitle}
                        onChange={(e) => handleInputChange("appTitle", e.target.value)}
                        placeholder="Your App Name"
                        className="mt-2"
                    />
                  </div>
                  <div className="space-y-3">
                    <Label htmlFor="appUrl" className="text-sm font-medium">
                      App URL
                    </Label>
                    <Input
                        id="appUrl"
                        name="appUrl"
                        type="url"
                        value={formData.appUrl}
                        onChange={(e) => handleInputChange("appUrl", e.target.value)}
                        placeholder="https://your-app.com"
                        className="mt-2"
                    />
                  </div>
                  <div className="space-y-3">
                    <Label htmlFor="jwtSecret" className="text-sm font-medium">
                      JWT Secret
                    </Label>
                    <Input
                        id="jwtSecret"
                        name="jwtSecret"
                        type="password"
                        value={formData.jwtSecret}
                        onChange={(e) => handleInputChange("jwtSecret", e.target.value)}
                        placeholder="Your JWT secret key"
                        className="mt-2"
                    />
                  </div>
                  <div className="space-y-3">
                    <Label htmlFor="maxProjectsPerOrg" className="text-sm font-medium">
                      Max Projects per Organization
                    </Label>
                    <Input
                        id="maxProjectsPerOrg"
                        name="maxProjectsPerOrg"
                        type="number"
                        value={formData.maxProjectsPerOrg}
                        onChange={(e) => handleInputChange("maxProjectsPerOrg", parseInt(e.target.value))}
                        min="1"
                        max="1000"
                        className="mt-2"
                    />
                  </div>
                </CardContent>
              </Card>

              {/* Feature Settings */}
              <Card className="h-fit">
                <CardHeader>
                  <div className="flex items-center gap-2">
                    <Shield className="h-5 w-5 text-green-500" />
                    <CardTitle>Feature Settings</CardTitle>
                  </div>
                  <CardDescription>
                    Control which features are enabled in your application
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="grid grid-cols-1 gap-6">
                    <div className="flex items-center justify-between py-2">
                      <Label htmlFor="allowRegistrations" className="text-sm font-medium">
                        Allow Registrations
                      </Label>
                      <Switch
                          id="allowRegistrations"
                          name="allowRegistrations"
                          checked={formData.allowRegistrations}
                          onCheckedChange={(checked) => handleInputChange("allowRegistrations", checked)}
                      />
                    </div>
                    <div className="flex items-center justify-between py-2">
                      <Label htmlFor="allowProjects" className="text-sm font-medium">
                        Allow Projects
                      </Label>
                      <Switch
                          id="allowProjects"
                          name="allowProjects"
                          checked={formData.allowProjects}
                          onCheckedChange={(checked) => handleInputChange("allowProjects", checked)}
                      />
                    </div>
                    <div className="flex items-center justify-between py-2">
                      <Label htmlFor="allowForms" className="text-sm font-medium">
                        Allow Forms
                      </Label>
                      <Switch
                          id="allowForms"
                          name="allowForms"
                          checked={formData.allowForms}
                          onCheckedChange={(checked) => handleInputChange("allowForms", checked)}
                      />
                    </div>
                    <div className="flex items-center justify-between py-2">
                      <Label htmlFor="allowStorage" className="text-sm font-medium">
                        Allow Storage
                      </Label>
                      <Switch
                          id="allowStorage"
                          name="allowStorage"
                          checked={formData.allowStorage}
                          onCheckedChange={(checked) => handleInputChange("allowStorage", checked)}
                      />
                    </div>
                    <div className="flex items-center justify-between py-2">
                      <Label htmlFor="allowBackups" className="text-sm font-medium">
                        Allow Backups
                      </Label>
                      <Switch
                          id="allowBackups"
                          name="allowBackups"
                          checked={formData.allowBackups}
                          onCheckedChange={(checked) => handleInputChange("allowBackups", checked)}
                      />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Storage Settings */}
            <Card>
              <CardHeader>
                <div className="flex items-center gap-2">
                  <Database className="h-5 w-5 text-purple-500" />
                  <CardTitle>Storage Settings</CardTitle>
                </div>
                <CardDescription>
                  Configure storage options, limits, and cloud provider settings
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-8">
                <div className="grid grid-cols-1 xl:grid-cols-2 gap-8">
                  <div className="space-y-6">
                    <div className="space-y-3">
                      <Label htmlFor="storageDriver" className="text-sm font-medium">
                        Storage Driver
                      </Label>
                      <Select
                          value={formData.storageDriver}
                          onValueChange={(value) => handleInputChange("storageDriver", value)}
                      >
                        <SelectTrigger className="mt-2">
                          <SelectValue placeholder="Select storage driver" />
                        </SelectTrigger>
                        <SelectContent>
                          {storageDriverOptions.map(option => (
                              <SelectItem key={option.value} value={option.value}>
                                {option.label}
                              </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-3">
                      <Label htmlFor="storageMaxContainers" className="text-sm font-medium">
                        Max Storage Containers
                      </Label>
                      <Input
                          id="storageMaxContainers"
                          name="storageMaxContainers"
                          type="number"
                          value={formData.storageMaxContainers}
                          onChange={(e) => handleInputChange("storageMaxContainers", parseInt(e.target.value))}
                          min="1"
                          className="mt-2"
                      />
                    </div>
                    <div className="space-y-3">
                      <Label htmlFor="storageMaxFileSizeInKB" className="text-sm font-medium">
                        Max File Size (KB)
                      </Label>
                      <Input
                          id="storageMaxFileSizeInKB"
                          name="storageMaxFileSizeInKB"
                          type="number"
                          value={formData.storageMaxFileSizeInKB}
                          onChange={(e) => handleInputChange("storageMaxFileSizeInKB", parseInt(e.target.value))}
                          min="1"
                          className="mt-2"
                      />
                    </div>
                  </div>
                  <div className="space-y-6">
                    <div className="space-y-3">
                      <Label htmlFor="storageAllowedMimes" className="text-sm font-medium">
                        Allowed MIME Types
                      </Label>
                      <Textarea
                          id="storageAllowedMimes"
                          name="storageAllowedMimes"
                          value={formData.storageAllowedMimes}
                          onChange={(e) => handleInputChange("storageAllowedMimes", e.target.value)}
                          placeholder="image/jpeg,image/png,application/pdf"
                          rows={4}
                          className="mt-2"
                      />
                    </div>
                  </div>
                </div>

                {/* Cloud Storage Configuration */}
                {formData.storageDriver !== "local" && (
                    <div className="pt-6 border-t">
                      {renderStorageConfiguration()}
                    </div>
                )}
              </CardContent>
            </Card>

            {/* API & Mail Settings */}
            <div className="grid grid-cols-1 xl:grid-cols-2 gap-8">
              {/* API Settings */}
              <Card className="h-fit">
                <CardHeader>
                  <div className="flex items-center gap-2">
                    <Shield className="h-5 w-5 text-orange-500" />
                    <CardTitle>API Settings</CardTitle>
                  </div>
                  <CardDescription>
                    Configure API throttling and rate limiting
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="flex items-center justify-between py-2">
                    <Label htmlFor="allowApiThrottle" className="text-sm font-medium">
                      Enable API Throttling
                    </Label>
                    <Switch
                        id="allowApiThrottle"
                        name="allowApiThrottle"
                        checked={formData.allowApiThrottle}
                        onCheckedChange={(checked) => handleInputChange("allowApiThrottle", checked)}
                    />
                  </div>
                  <div className="space-y-6 pt-4 border-t">
                    <div className="space-y-3">
                      <Label htmlFor="apiThrottleLimit" className="text-sm font-medium">
                        Throttle Limit
                      </Label>
                      <Input
                          id="apiThrottleLimit"
                          name="apiThrottleLimit"
                          type="number"
                          value={formData.apiThrottleLimit}
                          onChange={(e) => handleInputChange("apiThrottleLimit", parseInt(e.target.value))}
                          min="1"
                          className="mt-2"
                          disabled={!formData.allowApiThrottle}
                      />
                    </div>
                    <div className="space-y-3">
                      <Label htmlFor="apiThrottleInterval" className="text-sm font-medium">
                        Throttle Interval (seconds)
                      </Label>
                      <Input
                          id="apiThrottleInterval"
                          name="apiThrottleInterval"
                          type="number"
                          value={formData.apiThrottleInterval}
                          onChange={(e) => handleInputChange("apiThrottleInterval", parseInt(e.target.value))}
                          min="1"
                          className="mt-2"
                          disabled={!formData.allowApiThrottle}
                      />
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Mail Settings */}
              <Card className="h-fit">
                <CardHeader>
                  <div className="flex items-center gap-2">
                    <Mail className="h-5 w-5 text-blue-500" />
                    <CardTitle>Mail Settings</CardTitle>
                  </div>
                  <CardDescription>
                    Configure email service provider and settings
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="space-y-3">
                    <Label htmlFor="mailDriver" className="text-sm font-medium">
                      Mail Driver
                    </Label>
                    <Select
                        value={formData.mailDriver}
                        onValueChange={(value) => handleInputChange("mailDriver", value)}
                    >
                      <SelectTrigger className="mt-2">
                        <SelectValue placeholder="Select mail driver" />
                      </SelectTrigger>
                      <SelectContent>
                        {mailDriverOptions.map(option => (
                            <SelectItem key={option.value} value={option.value}>
                              {option.label}
                            </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>

                  {/* Mail Provider Configuration */}
                  <div className="pt-4 border-t">
                    {renderMailConfiguration()}
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Save Button */}
            <div className="flex justify-end pt-8 border-t bg-white p-6 rounded-lg shadow-sm">
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