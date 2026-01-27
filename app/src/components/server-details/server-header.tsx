import { ArrowLeft } from "@hugeicons/core-free-icons";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { HugeiconsIcon } from "@hugeicons/react";

interface ServerHeaderProps {
	name: string;
	serverType: string;
	serverVersion: string;
	status: string;
	onBack: () => void;
}

const getStatusColor = (status: string) => {
	switch (status.toLowerCase()) {
		case "running":
			return "default";
		case "stopped":
			return "secondary";
		case "starting":
			return "outline";
		case "stopping":
			return "destructive";
		default:
			return "secondary";
	}
};

export function ServerHeader({ name, serverType, serverVersion, status, onBack }: ServerHeaderProps) {
	return (
		<div className="flex items-center justify-between">
			<div className="flex items-center gap-4">
				<Button
					variant="ghost"
					size="icon"
					onClick={onBack}
				>
					<HugeiconsIcon icon={ArrowLeft} size={20} />
				</Button>
				<div>
					<h1 className="text-2xl font-bold">{name}</h1>
					<p className="text-muted-foreground text-sm mt-1">
						{serverType.charAt(0).toUpperCase() + serverType.slice(1)} {serverVersion}
					</p>
				</div>
			</div>
			<Badge variant={getStatusColor(status)}>
				{status}
			</Badge>
		</div>
	);
}
