import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import type { Route } from "./+types/page";
import { TrendingDownIcon, TrendingUpIcon } from "lucide-react";
import { Badge } from "~/components/ui/badge";
import { AppHeader } from "~/components/shared/header";
import { useQuery } from "@tanstack/react-query";
import { useEffect } from "react";

function getData(users, errors) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve({
        users,
        errors,
      });
    }, 1000);
  });
}

// export async function loader({ params }: Route.LoaderArgs) {
//   console.log("LOADER");
//   const promise = getData(200, 12);
//   return { promise };
// }

export async function clientLoader({ params }: Route.LoaderArgs) {
  const promise = getData(300, 15);
  return { promise };
}

export default function Dashboard({ loaderData }: Route.ComponentProps) {
  const { isLoading, isFetching, isError, data } = useQuery({
    queryKey: ["dashboard"],
    queryFn: async () => {
      const res = await loaderData.promise;
      return res;
    },
  });

  useEffect(() => {
    console.log(data);
  }, [data]);

  if (isLoading) {
    return <div>Loading</div>;
  }
  if (isFetching) {
    return <div>Fetching</div>;
  }

  const { users, errors } = data;
  return (
    <>
      <AppHeader title="Dashboard" />
      <div className="flex">
        <div className="*:data-[slot=card]:shadow-xs flex gap-4 py-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card lg:px-6">
          <Card className="">
            <CardHeader className="relative">
              <CardDescription>Total Users</CardDescription>
              <CardTitle className="@[250px]/card:text-3xl text-2xl font-semibold tabular-nums">
                {users}
              </CardTitle>
              <div className="absolute right-4 top-4">
                <Badge
                  variant="outline"
                  className="flex gap-1 rounded-lg text-xs"
                >
                  <TrendingUpIcon className="size-3" />
                  +12.5%
                </Badge>
              </div>
            </CardHeader>
            <CardFooter className="flex-col items-start gap-1 text-sm">
              <div className="line-clamp-1 flex gap-2 font-medium">
                Trending up this month <TrendingUpIcon className="size-4" />
              </div>
              <div className="text-muted-foreground">
                Visitors for the last 6 months
              </div>
            </CardFooter>
          </Card>
          <Card className="">
            <CardHeader className="relative">
              <CardDescription>New Errors</CardDescription>
              <CardTitle className="@[250px]/card:text-3xl text-2xl font-semibold tabular-nums">
                {errors}
              </CardTitle>
              <div className="absolute right-4 top-4">
                <Badge
                  variant="outline"
                  className="flex gap-1 rounded-lg text-xs"
                >
                  <TrendingDownIcon className="size-3" />
                  -20%
                </Badge>
              </div>
            </CardHeader>
            <CardFooter className="flex-col items-start gap-1 text-sm">
              <div className="line-clamp-1 flex gap-2 font-medium">
                Down 20% this period <TrendingDownIcon className="size-4" />
              </div>
              <div className="text-muted-foreground">
                Acquisition needs attention
              </div>
            </CardFooter>
          </Card>
        </div>
      </div>
    </>
  );
}
