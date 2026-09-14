import { cn } from "@/lib/utils";
import { DecorIcon } from "@/components/decor-icon";

type Logo = {
	src: string;
	alt: string;
	name: string;
};

export function LogoCloud() {
	return (
		<div className="grid grid-cols-2 border md:grid-cols-4">
			<LogoCard
				className="relative border-r border-b bg-secondary dark:bg-secondary/30"
				logo={{
					src: "https://cdn.simpleicons.org/go/black",
					alt: "Go Logo",
					name: "Go",
				}}
			>
				<DecorIcon className="z-10" position="bottom-right" />
			</LogoCard>

			<LogoCard
				className="border-b md:border-r"
				logo={{
					src: "https://cdn.simpleicons.org/gin/black",
					alt: "Gin Logo",
					name: "Gin",
				}}
			/>

			<LogoCard
				className="relative border-r border-b md:bg-secondary dark:md:bg-secondary/30"
				logo={{
					src: "https://cdn.simpleicons.org/postgresql/black",
					alt: "PostgreSQL Logo",
					name: "PostgreSQL",
				}}
			>
				<DecorIcon className="z-10" position="bottom-right" />
				<DecorIcon className="z-10 hidden md:block" position="bottom-left" />
			</LogoCard>

			<LogoCard
				className="relative border-b bg-secondary md:bg-background dark:bg-secondary/30 md:dark:bg-background"
				logo={{
					src: "https://cdn.simpleicons.org/jsonwebtokens/black",
					alt: "JWT Logo",
					name: "JWT",
				}}
			/>

			<LogoCard
				className="relative border-r border-b bg-secondary md:border-b-0 md:bg-background dark:bg-secondary/30 md:dark:bg-background"
				logo={{
					src: "https://cdn.simpleicons.org/react/black",
					alt: "React Logo",
					name: "React",
				}}
			>
				<DecorIcon className="z-10 md:hidden" position="bottom-right" />
			</LogoCard>

			<LogoCard
				className="border-b bg-background md:border-r md:border-b-0 md:bg-secondary dark:md:bg-secondary/30"
				logo={{
					src: "https://cdn.simpleicons.org/shadcnui/black",
					alt: "shadcn/ui Logo",
					name: "shadcn/ui",
				}}
			/>

			<LogoCard
				className="border-r"
				logo={{
					src: "https://cdn.simpleicons.org/tanstack/black",
					alt: "TanStack Query Logo",
					name: "TanStack Query",
				}}
			/>

			<LogoCard
				className="bg-secondary dark:bg-secondary/30"
				logo={{
					src: "https://cdn.simpleicons.org/typescript/black",
					alt: "TypeScript Logo",
					name: "TypeScript",
				}}
			/>
		</div>
	);
}

type LogoCardProps = React.ComponentProps<"div"> & {
	logo: Logo;
};

function LogoCard({ logo, className, children, ...props }: LogoCardProps) {
	return (
		<div
			className={cn(
				"flex items-center justify-center gap-2 bg-background px-4 py-8 md:p-8",
				className
			)}
			{...props}
		>
			<img
				alt={logo.alt}
				className="pointer-events-none h-4 shrink-0 select-none md:h-7 dark:brightness-0 dark:invert"
				height="auto"
				src={logo.src}
				width="auto"
			/>
			<span className="whitespace-nowrap font-bold text-foreground text-sm md:text-lg">
				{logo.name}
			</span>
			{children}
		</div>
	);
}