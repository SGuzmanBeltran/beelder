import { Card, CardContent } from "@/components/ui/card";

import { Button } from "@/components/ui/button";

export function Welcome() {
	return (
		<div className="flex items-center justify-center min-h-150">
			<Card className="max-w-2xl w-full">
				<CardContent className="flex flex-col items-center text-center py-16 px-8 space-y-6">
					{/* Minecraft-themed icon */}
					<div className="relative">
						<div className="size-24 bg-primary/10 rounded-lg flex items-center justify-center">
							<svg
								className="size-18 text-primary"
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 64 64"
								fill="none"
								stroke="currentColor"
								stroke-width="3"
								stroke-linecap="round"
								stroke-linejoin="round"
							>
								<rect x="16" y="10" width="32" height="10" rx="3" />
								<circle
									cx="44"
									cy="15"
									r="2"
									fill="currentColor"
									stroke="none"
								/>

								<rect x="12" y="24" width="40" height="12" rx="3" />
								<circle
									cx="46"
									cy="30"
									r="2"
									fill="currentColor"
									stroke="none"
								/>

								<rect x="8" y="40" width="48" height="14" rx="3" />
								<circle
									cx="48"
									cy="47"
									r="2"
									fill="currentColor"
									stroke="none"
								/>
							</svg>
						</div>
						<div className="absolute -top-1 -right-2 size-8 bg-primary rounded-md flex items-center justify-center p-1">
							<img src="/bee-icon.png" alt="Bee" className="object-contain" />
						</div>
					</div>

					{/* Welcome message */}
					<div className="space-y-2">
						<h2 className="text-2xl font-bold">Welcome to Beelder</h2>
						<p className="text-muted-foreground text-sm max-w-md">
							Create and manage Minecraft servers in minutes. Deploy Vanilla,
							Paper, Spigot, Forge servers and more with just a few clicks.
						</p>
					</div>

					{/* Features list */}
					<div className="grid grid-cols-1 md:grid-cols-3 gap-4 w-full max-w-xl text-left">
						<div className="space-y-1">
							<div className="flex items-center gap-2 text-sm font-medium">
								<div className="size-1.5 rounded-full bg-primary" />
								Easy Setup
							</div>
							<p className="text-xs text-muted-foreground pl-3.5">
								Automatic setup in minutes
							</p>
						</div>

						<div className="space-y-1">
							<div className="flex items-center gap-2 text-sm font-medium">
								<div className="size-1.5 rounded-full bg-primary" />
								Multiple Versions
							</div>
							<p className="text-xs text-muted-foreground pl-3.5">
								Vanilla, modded and more
							</p>
						</div>

						<div className="space-y-1">
							<div className="flex items-center gap-2 text-sm font-medium">
								<div className="size-1.5 rounded-full bg-primary" />
								Full Management
							</div>
							<p className="text-xs text-muted-foreground pl-3.5">
								Complete server control
							</p>
						</div>
					</div>

					{/* CTA Button */}
					<Button size="lg" className="mt-4" onClick={() => {}}>
						Create my first server
					</Button>

					<p className="text-xs text-muted-foreground">
						Start your Minecraft adventure now!
					</p>
				</CardContent>
			</Card>
		</div>
	);
}
