import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

interface ServerManagementProps {
	status: string;
}

export function ServerManagement({ status }: ServerManagementProps) {
	return (
		<Card>
			<CardHeader>
				<CardTitle>Server Management</CardTitle>
				<CardDescription>Control and manage your server</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				<div className="flex flex-wrap gap-3">
					<Button variant="default" disabled={status === "Running"}>
						Start Server
					</Button>
					<Button variant="destructive" disabled={status === "Stopped"}>
						Stop Server
					</Button>
					<Button variant="outline" disabled={status !== "Running"}>
						Restart Server
					</Button>
				</div>
				<Separator />
				<div className="flex flex-wrap gap-3">
					<Button variant="outline">
						Edit Configuration
					</Button>
					<Button variant="outline">
						View Console
					</Button>
					<Button variant="outline">
						File Manager
					</Button>
					<Button variant="outline">
						Manage Plugins
					</Button>
				</div>
			</CardContent>
		</Card>
	);
}
