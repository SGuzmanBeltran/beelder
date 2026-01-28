import { useLocation, useNavigate } from "react-router-dom";
import type { ServerDetails as ServerDetailsType } from "@/types/server";
import { ServerHeader } from "@/components/server-details/server-header";
import { ServerStatsGrid } from "@/components/server-details/server-stats-grid";
import { ServerManagement } from "@/components/server-details/server-management";
import { ServerConfiguration } from "@/components/server-details/server-configuration";
import { ConnectionDetails } from "@/components/server-details/connection-details";
import { ResourceUsage } from "@/components/server-details/resource-usage";

export function ServerDetails() {
	const location = useLocation();
	const navigate = useNavigate();

	// Get server data from location state
	const server: ServerDetailsType = location.state?.server || {
		id: "srv_a1b2c3d4e5f6",
		name: "My Awesome Server",
		server_version: "1.20.1",
		server_type: "paper",
		region: "us-east-1",
		player_count: 20,
		ram_plan: "4GB",
		difficulty: "normal",
		online_mode: true,
		status: "Running",
		ip_address: "mc.example.com",
		port: 25565,
		created_at: new Date().toISOString(),
		current_players: 5,
		cpu_usage: 45,
		ram_usage: 68,
		disk_usage: 2.4,
		network_in: 1.2,
		network_out: 0.8,
		uptime: "3d 12h 45m",
		tps: 19.8,
		world_size: 1.8,
		motd: "§6Welcome to §bMy Awesome Server§r\\n§7Join the adventure!",
		world_seed: "1234567890",
		last_backup: "2 hours ago",
		gamemode: "survival",
		pvp_enabled: true,
	};

	return (
		<div className="min-h-screen px-6 lg:px-12 pt-8">
			<div className="max-w-7xl mx-auto space-y-8">
				{/* Header */}
				<ServerHeader
					name={server.name}
					serverType={server.server_type}
					serverVersion={server.server_version}
					status={server.status}
					onBack={() => navigate("/")}
				/>

				{/* Quick Stats Grid */}
				<ServerStatsGrid
					currentPlayers={server.current_players || 0}
					maxPlayers={server.player_count}
					tps={server.tps || 20.0}
					uptime={server.uptime || "0h 0m"}
					worldSize={server.world_size || 0.0}
				/>

				{/* Main Content Grid */}
				<div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
					{/* Left Column - Server Info & Connection */}
					<div className="lg:col-span-2 space-y-6">
						<ServerManagement status={server.status} />

						<ServerConfiguration
							serverId={server.id}
							region={server.region}
							ramPlan={server.ram_plan}
							difficulty={server.difficulty}
							gamemode={server.gamemode || "survival"}
							pvpEnabled={server.pvp_enabled !== false}
							onlineMode={server.online_mode}
							worldSeed={server.world_seed || ""}
							createdAt={server.created_at}
						/>

						<ConnectionDetails
							ipAddress={server.ip_address}
							port={server.port}
						/>
					</div>

					{/* Right Column - Performance & Stats */}
					<div className="space-y-6">
						<ResourceUsage
							cpuUsage={server.cpu_usage || 0}
							ramUsage={server.ram_usage || 0}
							diskUsage={server.disk_usage || 0}
							networkIn={server.network_in || 0}
							networkOut={server.network_out || 0}
						/>
					</div>
				</div>
			</div>
		</div>
	);
}