import { Card, CardContent } from "./ui/card";

import { Button } from "./ui/button";

interface PricingCardProps {
	ram: string;
	price: string;
	badge?: {
		text: string;
		color: "red" | "green" | "yellow" | "purple" | "stone";
	};
	disabled?: boolean;
	loading?: boolean;
}

const badgeColors = {
	red: "bg-red-900 border-red-900/50",
	green: "bg-green-600 border-green-500",
	yellow: "bg-yellow-500 border-yellow-400",
	purple: "bg-purple-600 border-purple-500",
	stone: "bg-stone-700 border-stone-700",
};

const cardBorderColors = {
	red: "border-red-900/50",
	green: "border-green-500",
	yellow: "border-yellow-400",
	purple: "border-purple-500",
	stone: "border-stone-700",
};

export function PricingCard({
	ram,
	price,
	badge,
	disabled = false,
	loading = false,
}: PricingCardProps) {
	const borderColor = badge
		? cardBorderColors[badge.color]
		: "border-stone-700";

	return (
		<Card
			className={`relative bg-stone-900 border-2 ${borderColor} overflow-visible min-w-70 shrink-0`}
		>
			<CardContent className="p-6 space-y-4 overflow-visible">
				{badge && (
					<div className="absolute -top-4 left-1/2 -translate-x-1/2 z-10">
						<div
							className={`${
								badgeColors[badge.color]
							} text-white px-4 py-1 rounded-full text-sm font-semibold border-2 whitespace-nowrap`}
						>
							{badge.text}
						</div>
					</div>
				)}
				<div className="text-center pt-4">
					<h3 className="text-4xl font-bold text-white">{ram}</h3>
				</div>
				<div className="text-center">
					<p className="text-3xl font-bold text-white">{price}</p>
					<p className="text-sm text-muted-foreground">/ MONTH</p>
				</div>
				<Button
					type="submit"
					className={`w-full text-black ${
						disabled
							? "bg-stone-700 hover:bg-stone-600"
							: "bg-yellow-500 hover:bg-yellow-700"
					}`}
					disabled={disabled}
				>
					{loading ? "Creating..." : "Create server"}
				</Button>
			</CardContent>
		</Card>
	);
}
