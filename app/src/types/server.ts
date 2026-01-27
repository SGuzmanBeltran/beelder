export interface ServerDetails {
	id: string;
	name: string;
	server_version: string;
	server_type: string;
	region: string;
	player_count: number;
	ram_plan: string;
	difficulty: string;
	online_mode: boolean;
	status: string; // e.g., "Running", "Stopped", "Starting", "Stopping"
	ip_address: string;
	port: number;
	created_at: string;
	// Additional metrics
	current_players?: number;
	cpu_usage?: number; // percentage
	ram_usage?: number; // percentage
	disk_usage?: number; // in GB
	network_in?: number; // in MB/s
	network_out?: number; // in MB/s
	uptime?: string;
	tps?: number; // ticks per second
	world_size?: number; // in GB
	motd?: string;
	world_seed?: string;
	last_backup?: string;
	gamemode?: string;
	pvp_enabled?: boolean;
}
