import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

interface ResourceUsageProps {
	cpuUsage: number;
	ramUsage: number;
	diskUsage: number;
	networkIn: number;
	networkOut: number;
}

export function ResourceUsage({ cpuUsage, ramUsage, diskUsage, networkIn, networkOut }: ResourceUsageProps) {
	return (
		<Card>
			<CardHeader>
				<CardTitle>Resource Usage</CardTitle>
				<CardDescription>Real-time server metrics</CardDescription>
			</CardHeader>
			<CardContent className="space-y-6">
				<div className="space-y-2">
					<div className="flex items-center justify-between text-sm">
						<span className="text-muted-foreground">CPU Usage</span>
						<span className="font-semibold">{cpuUsage}%</span>
					</div>
					<div className="w-full bg-muted rounded-full h-2">
						<div
							className="bg-primary h-2 rounded-full transition-all"
							style={{ width: `${cpuUsage}%` }}
						/>
					</div>
				</div>
				<div className="space-y-2">
					<div className="flex items-center justify-between text-sm">
						<span className="text-muted-foreground">RAM Usage</span>
						<span className="font-semibold">{ramUsage}%</span>
					</div>
					<div className="w-full bg-muted rounded-full h-2">
						<div
							className="bg-primary h-2 rounded-full transition-all"
							style={{ width: `${ramUsage}%` }}
						/>
					</div>
				</div>
				<Separator />
				<div className="space-y-3">
					<div className="flex items-center justify-between text-sm">
						<span className="text-muted-foreground">Disk Usage</span>
						<span className="font-semibold">{diskUsage.toFixed(1)} GB</span>
					</div>
					<div className="flex items-center justify-between text-sm">
						<span className="text-muted-foreground">Network In</span>
						<span className="font-semibold">{networkIn.toFixed(1)} MB/s</span>
					</div>
					<div className="flex items-center justify-between text-sm">
						<span className="text-muted-foreground">Network Out</span>
						<span className="font-semibold">{networkOut.toFixed(1)} MB/s</span>
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
