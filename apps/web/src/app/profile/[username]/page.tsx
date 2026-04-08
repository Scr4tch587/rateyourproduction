import { notFound } from "next/navigation";
import { Separator } from "@/components/ui/separator";
import { apiFetch } from "@/lib/api";
import { LogEntryItem } from "@/components/log-entry";
import type { Profile, PaginatedResponse, LogEntry } from "@/lib/types";

async function getProfile(username: string) {
  try {
    return await apiFetch<Profile>(`/api/v1/profile/${username}`);
  } catch {
    return null;
  }
}

async function getUserLogs(userId: string) {
  try {
    return await apiFetch<PaginatedResponse<LogEntry>>(`/api/v1/logs?user_id=${userId}&per_page=50`);
  } catch {
    return null;
  }
}

interface ProfilePageProps {
  params: Promise<{ username: string }>;
}

export default async function ProfilePage({ params }: ProfilePageProps) {
  const { username } = await params;
  const profile = await getProfile(username);

  if (!profile) {
    notFound();
  }

  const logs = await getUserLogs(profile.id);

  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <div className="mb-4">
        <h1 className="font-serif text-xl font-bold">
          {profile.display_name || profile.username}
        </h1>
        <div className="text-xs text-muted-foreground">
          @{profile.username} &middot; {profile.log_count} log{profile.log_count !== 1 ? "s" : ""} &middot; {profile.review_count} review{profile.review_count !== 1 ? "s" : ""}
        </div>
      </div>

      <Separator />

      <div className="mt-6">
        <div className="border-b border-border pb-1 mb-3">
          <h2 className="font-serif text-sm font-bold">Theatre Journal</h2>
        </div>
        {logs?.data && logs.data.length > 0 ? (
          <div>
            {logs.data.map((log) => (
              <LogEntryItem key={log.id} log={log} showUser={false} />
            ))}
          </div>
        ) : (
          <div className="text-xs text-muted-foreground py-6 italic">
            No entries yet.
          </div>
        )}
      </div>
    </div>
  );
}
