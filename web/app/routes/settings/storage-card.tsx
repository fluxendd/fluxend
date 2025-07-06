// StorageSettingsCard.tsx
import { Database, Cloud } from "lucide-react";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { Switch } from "~/components/ui/switch";
import { Textarea } from "~/components/ui/textarea";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "~/components/ui/card";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "~/components/ui/select";
import type { SettingsData } from "~/services/settings";

interface StorageSettingsCardProps {
    formData: SettingsData;
    onInputChange: (field: keyof SettingsData, value: any) => void;
}

const storageDriverOptions = [
    { value: "local", label: "Local Storage" },
    { value: "aws", label: "AWS S3" },
    { value: "backblaze", label: "Backblaze B2" },
    { value: "dropbox", label: "Dropbox" },
];

export const StorageSettingsCard = ({ formData, onInputChange }: StorageSettingsCardProps) => {
    const renderStorageConfiguration = () => {
        const isStorageDisabled = !formData.allowStorage;

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
                                    onChange={(e) => onInputChange("awsAccessKeyId", e.target.value)}
                                    placeholder="Your AWS Access Key ID"
                                    className="mt-2"
                                    disabled={isStorageDisabled}
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
                                    onChange={(e) => onInputChange("awsSecretAccessKey", e.target.value)}
                                    placeholder="Your AWS Secret Access Key"
                                    className="mt-2"
                                    disabled={isStorageDisabled}
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
                                    onChange={(e) => onInputChange("awsRegion", e.target.value)}
                                    placeholder="us-east-1"
                                    className="mt-2"
                                    disabled={isStorageDisabled}
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
                                    onChange={(e) => onInputChange("backblazeKeyId", e.target.value)}
                                    placeholder="Your Backblaze Key ID"
                                    className="mt-2"
                                    disabled={isStorageDisabled}
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
                                    onChange={(e) => onInputChange("backblazeApplicationKey", e.target.value)}
                                    placeholder="Your Backblaze Application Key"
                                    className="mt-2"
                                    disabled={isStorageDisabled}
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
                                    onChange={(e) => onInputChange("dropboxAccessToken", e.target.value)}
                                    placeholder="Your Dropbox Access Token"
                                    className="mt-2"
                                    disabled={isStorageDisabled}
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
                                    onChange={(e) => onInputChange("dropboxAppKey", e.target.value)}
                                    placeholder="Your Dropbox App Key"
                                    className="mt-2"
                                    disabled={isStorageDisabled}
                                />
                            </div>
                        </div>
                    </div>
                );
            default:
                return null;
        }
    };

    return (
        <Card className="h-fit">
            <CardHeader>
                <div className="flex items-center gap-2">
                    <Database className="h-5 w-5 text-purple-500" />
                    <CardTitle>Storage Settings</CardTitle>
                </div>
                <CardDescription>
                    Configure storage options, limits, and cloud provider settings
                </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
                <div className="flex items-center justify-between py-2">
                    <Label htmlFor="allowStorage" className="text-sm font-medium">
                        Allow Storage
                    </Label>
                    <Switch
                        id="allowStorage"
                        name="allowStorage"
                        checked={formData.allowStorage}
                        onCheckedChange={(checked) => onInputChange("allowStorage", checked)}
                    />
                </div>

                <div className="space-y-6 pt-4 border-t">
                    <div className="space-y-3">
                        <Label htmlFor="storageDriver" className="text-sm font-medium">
                            Storage Driver
                        </Label>
                        <Select
                            value={formData.storageDriver}
                            onValueChange={(value) => onInputChange("storageDriver", value)}
                            disabled={!formData.allowStorage}
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
                            onChange={(e) => onInputChange("storageMaxContainers", parseInt(e.target.value))}
                            min="1"
                            className="mt-2"
                            disabled={!formData.allowStorage}
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
                            onChange={(e) => onInputChange("storageMaxFileSizeInKB", parseInt(e.target.value))}
                            min="1"
                            className="mt-2"
                            disabled={!formData.allowStorage}
                        />
                    </div>
                    <div className="space-y-3">
                        <Label htmlFor="storageAllowedMimes" className="text-sm font-medium">
                            Allowed MIME Types
                        </Label>
                        <Textarea
                            id="storageAllowedMimes"
                            name="storageAllowedMimes"
                            value={formData.storageAllowedMimes}
                            onChange={(e) => onInputChange("storageAllowedMimes", e.target.value)}
                            placeholder="image/jpeg,image/png,application/pdf"
                            rows={4}
                            className="mt-2"
                            disabled={!formData.allowStorage}
                        />
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
    );
};