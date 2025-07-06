// MailSettingsCard.tsx
import { Mail } from "lucide-react";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
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

interface MailSettingsCardProps {
    formData: SettingsData;
    onInputChange: (field: keyof SettingsData, value: any) => void;
}

const mailDriverOptions = [
    { value: "sendgrid", label: "SendGrid" },
    { value: "mailgun", label: "Mailgun" },
];

const mailgunRegionOptions = [
    { value: "us", label: "United States" },
    { value: "eu", label: "Europe" },
];

export const MailSettingsCard = ({ formData, onInputChange }: MailSettingsCardProps) => {
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
                                    onChange={(e) => onInputChange("sendgridApiKey", e.target.value)}
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
                                    onChange={(e) => onInputChange("sendgridEmailSource", e.target.value)}
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
                                    onChange={(e) => onInputChange("mailgunApiKey", e.target.value)}
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
                                    onChange={(e) => onInputChange("mailgunEmailSource", e.target.value)}
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
                                    onChange={(e) => onInputChange("mailgunDomain", e.target.value)}
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
                                    onValueChange={(value) => onInputChange("mailgunRegion", value)}
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
                        onValueChange={(value) => onInputChange("mailDriver", value)}
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
    );
};