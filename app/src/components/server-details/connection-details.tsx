import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { HugeiconsIcon } from "@hugeicons/react";
import { Copy } from "@hugeicons/core-free-icons";

interface ConnectionDetailsProps {
	ipAddress: string;
	port: number;
}

export function ConnectionDetails({ ipAddress, port }: ConnectionDetailsProps) {
	return (
		<Card>
			<CardHeader>
				<CardTitle>Connection Details</CardTitle>
				<CardDescription>Use these details to connect to your server</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				<div className="grid grid-cols-1 md:grid-cols-2 gap-6">
					<div className="space-y-2">
						<h3 className="text-sm font-medium text-muted-foreground">Server Address</h3>
						<div className="flex items-center gap-2">
							<code className="flex-1 px-3 py-2 bg-muted rounded-md text-sm font-mono">
								{ipAddress}
							</code>
							<Button variant="outline" size="icon">
								<HugeiconsIcon icon={Copy} size={16} />
							</Button>
						</div>
					</div>
					<div className="space-y-2">
						<h3 className="text-sm font-medium text-muted-foreground">Port</h3>
						<div className="flex items-center gap-2">
							<code className="flex-1 px-3 py-2 bg-muted rounded-md text-sm font-mono">
								{port}
							</code>
							<Button variant="outline" size="icon">
								<HugeiconsIcon icon={Copy} size={16} />
							</Button>
						</div>
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
