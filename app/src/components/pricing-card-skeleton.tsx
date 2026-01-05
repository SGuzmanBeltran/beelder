import { Card, CardContent } from "./ui/card";

export function PricingCardSkeleton() {
	return (
		<div className="relative flex items-center justify-center gap-4">
			{/* Invisible Previous Button Placeholder */}
			<div className="w-10 h-10 shrink-0"></div>

			{/* Skeleton Card */}
			<div className="flex-1 animate-pulse">
				<Card className="border-2 border-stone-700 bg-stone-900 min-w-70 shrink-0">
					<CardContent className="p-6 space-y-4">
						{/* RAM size skeleton */}
						<div className="text-center pt-3">
							<div className="h-9 bg-stone-700 rounded w-24 mx-auto"></div>
						</div>
						{/* Price skeleton */}
						<div className="text-center space-y-2">
							<div className="h-8 bg-stone-700 rounded w-28 mx-auto"></div>
							<div className="h-3 bg-stone-700 rounded w-20 mx-auto"></div>
						</div>
						{/* Button skeleton */}
						<div className="h-10 bg-stone-700 rounded w-full"></div>
					</CardContent>
				</Card>
			</div>

			{/* Invisible Next Button Placeholder */}
			<div className="w-10 h-10 shrink-0"></div>
		</div>
	);
}
