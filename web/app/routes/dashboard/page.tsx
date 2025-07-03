import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import type { Route } from "./+types/page";
import {
  CheckCircleIcon,
  XCircleIcon,
  Database,
  Server,
  HardDrive,
  Cpu,
} from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";
import { AppHeader } from "~/components/shared/header";
import { redirect, useOutletContext } from "react-router";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import { getAuthToken, getClientAuthToken } from "~/lib/auth";
import { initializeServices } from "~/services";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Dashboard - Fluxend" },
    { name: "description", content: "Fluxend Dashboard" },
  ];
}

function StatusCard({
  title,
  status,
  icon: Icon,
  isStale,
}: {
  title: string;
  status: string;
  icon: React.ComponentType<{ className?: string }>;
  isStale: boolean;
}) {
  return (
    <Card className="flex-1 min-w-[120px] rounded-lg">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <Icon className="size-4 text-muted-foreground" />
          {isStale ? (
            <XCircleIcon className="size-4 text-red-500" />
          ) : (
            <CheckCircleIcon className="size-4 text-green-500" />
          )}
        </div>
        <CardDescription className="text-xs">{title}</CardDescription>
        <CardTitle className="text-lg font-semibold">{status}</CardTitle>
      </CardHeader>
    </Card>
  );
}

function MetricCard({
  title,
  value,
  subtitle,
  icon: Icon,
  data,
  isStale,
}: {
  title: string;
  value: string;
  subtitle?: string;
  icon: React.ComponentType<{ className?: string }>;
  data: number[];
  isStale: boolean;
}) {
  const maxValue = Math.max(...data);
  const minValue = Math.min(...data);

  return (
    <Card className="flex-1 min-w-[200px] rounded-lg">
      <CardHeader className="relative pb-2">
        <div className="flex items-center justify-between">
          <Icon className="size-4 text-muted-foreground" />
          {isStale ? (
            <XCircleIcon className="size-4 text-red-500" />
          ) : (
            <span className="relative flex size-3">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-green-500 opacity-75"></span>
              <span className="relative inline-flex size-3 rounded-full bg-green-500"></span>
            </span>
          )}
        </div>
        <CardDescription className="text-xs">{title}</CardDescription>
        <CardTitle className="text-2xl font-semibold tabular-nums">
          {value}
        </CardTitle>
      </CardHeader>
      <CardFooter className="pt-0">
        <div className="w-full">
          {subtitle && (
            <div className="text-xs text-muted-foreground">{subtitle}</div>
          )}
        </div>
      </CardFooter>
    </Card>
  );
}

export async function clientLoader() {
  const authToken = await getClientAuthToken();
  if (!authToken) {
    return redirect("/logout");
  }

  const services = initializeServices(authToken);

  const data = await services.dashboard.getHealthStatus();

  return data;
}

export default function Dashboard({ loaderData }: Route.ComponentProps) {
  const { services } = useOutletContext<ProjectLayoutOutletContext>();

  const {
    isLoading,
    isError,
    data: healthData,
    isStale,
  } = useQuery({
    queryKey: ["dashboard-health"],
    initialData: loaderData,
    queryFn: async () => {
      return await services.dashboard.getHealthStatus();
    },
    staleTime: 15000,
    refetchInterval: 10000, // Refetch every 10 seconds
    refetchIntervalInBackground: true,
  });

  // Generate mock historical data for charts
  const diskData = useMemo(() => {
    const baseValue = parseFloat(
      healthData?.disk_usage?.replace("%", "") || "48.9"
    );
    return Array.from({ length: 8 }, (_, i) =>
      Math.max(0, Math.min(100, baseValue + (Math.random() - 0.5) * 10))
    );
  }, [healthData?.disk_usage]);

  const cpuData = useMemo(() => {
    const baseValue = parseFloat(
      healthData?.cpu_usage?.replace("%", "") || "0"
    );
    return Array.from({ length: 8 }, (_, i) =>
      Math.max(0, Math.min(100, baseValue + (Math.random() - 0.5) * 5))
    );
  }, [healthData?.cpu_usage]);

  if (isLoading) {
    return (
      <>
        <AppHeader title="Dashboard" />
        <div className="flex flex-col gap-4 p-4">
          <div className="text-center">Loading dashboard...</div>
        </div>
      </>
    );
  }

  if (isError || !healthData) {
    return (
      <>
        <AppHeader title="Dashboard" isLoading={false} loadingProgress={0} />
        <div className="flex flex-col gap-4 p-4">
          <div className="text-center text-red-500">
            Error loading dashboard data
          </div>
        </div>
      </>
    );
  }
  return (
    <>
      <AppHeader title="Dashboard" />
      <div className="flex flex-col gap-6 p-4">
        <div className="flex flex-wrap gap-4">
          <StatusCard
            title="Database"
            status={healthData.database_status}
            icon={Database}
            isStale={!isLoading && isStale}
          />
          <StatusCard
            title="UI App"
            status={healthData.app_status}
            icon={Server}
            isStale={!isLoading && isStale}
          />
          <StatusCard
            title="PostgREST"
            status={healthData.postgrest_status}
            icon={Server}
            isStale={!isLoading && isStale}
          />
        </div>

        {/* Metrics Cards Row */}
        <div className="flex flex-wrap gap-4">
          <MetricCard
            title="Disk Usage"
            value={healthData.disk_usage}
            subtitle={`${healthData.disk_available} available of ${healthData.disk_total}`}
            icon={HardDrive}
            data={diskData}
            isStale={isStale}
          />
          <MetricCard
            title="CPU Usage"
            value={healthData.cpu_usage}
            subtitle={`${healthData.cpu_cores} core${
              healthData.cpu_cores !== 1 ? "s" : ""
            } • Real-time utilization`}
            icon={Cpu}
            data={cpuData}
            isStale={isStale}
          />
        </div>
      </div>
    </>
  );
}
