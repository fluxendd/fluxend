// ApiSettingsCard.tsx
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

interface ApiSettingsCardProps {
    formData: SettingsData;
    onInputChange: (field: keyof SettingsData, value: any) => void;
}

export const ApiSettingsCard = ({ formData, onInputChange }: ApiSettingsCardProps) => {
    return (
        <Card className="h-fit">
            <CardHeader>
                <div className="flex items-center gap-2">
                    <Settings className="h-5 w-5 text-orange-500" />
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
                        onCheckedChange={(checked) => onInputChange("allowApiThrottle", checked)}
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
                            onChange={(e) => onInputChange("apiThrottleLimit", parseInt(e.target.value))}
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
                            onChange={(e) => onInputChange("apiThrottleInterval", parseInt(e.target.value))}
                            min="1"
                            className="mt-2"
                            disabled={!formData.allowApiThrottle}
                        />
                    </div>
                </div>
            </CardContent>
        </Card>
    );
};