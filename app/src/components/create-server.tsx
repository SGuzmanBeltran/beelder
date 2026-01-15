import { ArrowLeft, ArrowRight } from "@hugeicons/core-free-icons";
import { Card, CardContent } from "./ui/card";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "./ui/select";
import { pricingPlans, useServerCreation } from "@/hooks/useServerCreation";

import { HugeiconsIcon } from "@hugeicons/react";
import { Input } from "./ui/input";
import { PricingCard } from "./pricing-card";
import { PricingCardSkeleton } from "./pricing-card-skeleton";
import { Slider } from "./ui/slider";
import { Switch } from "./ui/switch";

export function CreateServer() {
	const {
		currentPlanIndex,
		isLoadingRecommendation,
		recommendedPlan,
		isSubmitting,
		form,
		handlePrevious,
		handleNext,
		sendFormData,
		setCurrentPlanIndex,
		currentPlan,
	} = useServerCreation();

	return (
		<form
			onSubmit={form.handleSubmit(sendFormData)}
			className="flex flex-col items-center justify-center min-h-150 space-y-8 px-4 w-full lg:w-2/3 lg:px-0"
		>
			<div className="flex w-full justify-start">
				<h1 className=" text-2xl font-bold">Create your server</h1>
			</div>

			<div className="flex flex-col sm:flex-row w-full space-y-12 sm:space-x-8">
				<div className="flex flex-col w-full sm:w-4/5 space-y-8">
					<Card className="w-full">
						<CardContent className="space-y-5">
							<h3 className="text-lg">Configure your installation</h3>
							<div className="space-y-3">
								<h4>What should we install on you server?</h4>
								<Select
									onValueChange={(value) => form.setValue("serverType", value)}
								>
									<SelectTrigger className="w-full my-2">
										<SelectValue placeholder="Select server type" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="vanilla">Vanilla</SelectItem>
										<SelectItem value="paper">Paper</SelectItem>
										<SelectItem value="forge">Forge</SelectItem>
									</SelectContent>
								</Select>

								{form.formState.errors.serverType && (
									<p className="text-red-500">
										{form.formState.errors.serverType.message as string}
									</p>
								)}

								<Select
									onValueChange={(value) =>
										form.setValue("serverVersion", value)
									}
								>
									<SelectTrigger className="w-full my-2">
										<SelectValue placeholder="Select server version" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="1.21.11">1.21.11</SelectItem>
										<SelectItem value="1.21.10">1.21.10</SelectItem>
										<SelectItem value="1.21.9">1.21.9</SelectItem>
									</SelectContent>
								</Select>
								{form.formState.errors.serverVersion && (
									<p className="text-red-500">
										{form.formState.errors.serverVersion.message as string}
									</p>
								)}
							</div>
							<div className="space-y-3">
								<h4>
									How many players do you expect to be on the server at once?
								</h4>
								<Slider
									className="my-2"
									value={
										form.watch("playerCount")
											? [form.watch("playerCount")]
											: [0]
									}
									onValueChange={(value) =>
										form.setValue("playerCount", value[0])
									}
									defaultValue={[0]}
									max={100}
									step={1}
								/>
								<p className="text-sm text-stone-400">
									{form.watch("playerCount")} players
								</p>
								{form.formState.errors.playerCount && (
									<p className="text-red-500">
										{form.formState.errors.playerCount.message as string}
									</p>
								)}
							</div>
						</CardContent>
					</Card>

					<Card className="w-full">
						<CardContent className="space-y-3">
							<h3 className="text-lg">Select your location</h3>
							<Select onValueChange={(value) => form.setValue("region", value)}>
								<SelectTrigger className="w-full my-2">
									<SelectValue placeholder="Select a region" />
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="us-east-1">
										US East (N. Virginia)
									</SelectItem>
									<SelectItem value="us-west-2">US West (Oregon)</SelectItem>
									<SelectItem value="eu-west-1">EU West (Ireland)</SelectItem>
									<SelectItem value="ap-southeast-1">
										Asia Pacific (Singapore)
									</SelectItem>
									<SelectItem value="sa-east-1">
										South America (São Paulo)
									</SelectItem>
								</SelectContent>
							</Select>
							{form.formState.errors.region && (
								<p className="text-red-500">
									{form.formState.errors.region.message as string}
								</p>
							)}
						</CardContent>
					</Card>

					<Card className="w-full pb-8">
						<CardContent className="space-y-4">
							<h3 className="text-lg">Add your initial configuration</h3>
							<div className="space-y-3">
								<h4>Server name</h4>
								<Input
									onChange={(e) => form.setValue("serverName", e.target.value)}
								/>
								{form.formState.errors.serverName && (
									<p className="text-red-500">
										{form.formState.errors.serverName.message as string}
									</p>
								)}
							</div>

							<div className="flex justify-between space-x-6">
								<div className="w-1/2 space-y-3">
									<h4>Select difficulty</h4>
									<Select
										onValueChange={(value) =>
											form.setValue("difficulty", value)
										}
									>
										<SelectTrigger className="w-full">
											<SelectValue placeholder="Select a difficulty" />
										</SelectTrigger>
										<SelectContent>
											<SelectItem value="peaceful">Peaceful</SelectItem>
											<SelectItem value="easy">Easy</SelectItem>
											<SelectItem value="normal">Normal</SelectItem>
											<SelectItem value="hard">Hard</SelectItem>
											<SelectItem value="hardcore">Hardcore</SelectItem>
										</SelectContent>
									</Select>
									{form.formState.errors.difficulty && (
										<p className="text-red-500">
											{form.formState.errors.difficulty.message as string}
										</p>
									)}
								</div>
								<div className="w-1/2 flex flex-col justify-center space-y-3 md:space-y-1">
									<h4 className="pb-2">Allow only premium players</h4>
									<Switch
										checked={form.watch("onlineMode") || false}
										onCheckedChange={(checked) =>
											form.setValue("onlineMode", checked)
										}
									/>
								</div>
							</div>
						</CardContent>
					</Card>
				</div>

				<div className="flex flex-col w-full sm:w-2/5 space-y-8">
					<div className="max-w-2xl w-full overflow-visible">
						<div className="space-y-4 py-8">
							{isLoadingRecommendation || recommendedPlan === null ? (
								<PricingCardSkeleton />
							) : (
								<>
									{/* Carousel */}
									<div className="relative flex items-center justify-center gap-4">
										{/* Previous Button */}
										<button
											type="button"
											onClick={handlePrevious}
											className="z-10 flex items-center justify-center w-10 h-10 rounded-full bg-stone-800 hover:bg-stone-700 border-2 border-stone-600 transition-colors shrink-0"
											aria-label="Previous plan"
										>
											<HugeiconsIcon icon={ArrowLeft} size={20} />
										</button>

										{/* Pricing Card */}
										<div className="flex-1">
											<PricingCard
												ram={currentPlan.ram}
												price={currentPlan.price}
												badge={currentPlan.badge}
												loading={isSubmitting}
												disabled={isSubmitting}
											/>
										</div>

										{/* Next Button */}
										<button
											type="button"
											onClick={handleNext}
											className="z-10 flex items-center justify-center w-10 h-10 rounded-full bg-stone-800 hover:bg-stone-700 border-2 border-stone-600 transition-colors shrink-0"
											aria-label="Next plan"
										>
											<HugeiconsIcon icon={ArrowRight} size={20} />
										</button>
									</div>

									{/* Indicator Dots */}
									<div className="flex justify-center gap-2">
										{pricingPlans.map((_, index) => (
											<button
												key={index}
												type="button"
												onClick={() => {
													setCurrentPlanIndex(index);
													form.setValue("ramPlan", pricingPlans[index].ram);
												}}
												className={`w-2 h-2 rounded-full transition-all ${
													index === currentPlanIndex
														? "bg-yellow-500 w-8"
														: "bg-stone-600 hover:bg-stone-500"
												}`}
												aria-label={`Go to plan ${index + 1}`}
											/>
										))}
									</div>
								</>
							)}
						</div>
					</div>
				</div>
			</div>
		</form>
	);
}
