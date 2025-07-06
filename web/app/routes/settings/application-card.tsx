// ApplicationSettingsCard.tsx
import { Settings } from "lucide-react";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { Switch } from "~/components/ui/switch";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "~/components/ui/card";
import type { SettingsData } from "~/services/settings";

interface ApplicationSettingsCardProps {
    formData: SettingsData;
    onInputChange: (field: keyof SettingsData, value: any) => void;
}

export const ApplicationSettingsCard = ({ formData, onInputChange }: ApplicationSettingsCardProps) => {
    return (
        <Card className="h-fit">
            <CardHeader>
                <div className="flex items-center gap-2">
                    <Settings className="h-5 w-5 text-amber-400" />
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
                        onChange={(e) => onInputChange("appTitle", e.target.value)}
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
                        onChange={(e) => onInputChange("appUrl", e.target.value)}
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
                        onChange={(e) => onInputChange("jwtSecret", e.target.value)}
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
                        onChange={(e) => onInputChange("maxProjectsPerOrg", parseInt(e.target.value))}
                        min="1"
                        max="1000"
                        className="mt-2"
                    />
                </div>
                <div className="pt-4 border-t">
                    <div className="flex items-center justify-between py-2">
                        <Label htmlFor="allowRegistrations" className="text-sm font-medium">
                            Allow Registrations
                        </Label>
                        <Switch
                            id="allowRegistrations"
                            name="allowRegistrations"
                            checked={formData.allowRegistrations}
                            onCheckedChange={(checked) => onInputChange("allowRegistrations", checked)}
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
                            onCheckedChange={(checked) => onInputChange("allowProjects", checked)}
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
                            onCheckedChange={(checked) => onInputChange("allowForms", checked)}
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
                            onCheckedChange={(checked) => onInputChange("allowBackups", checked)}
                        />
                    </div>
                </div>
            </CardContent>
        </Card>
    );
};