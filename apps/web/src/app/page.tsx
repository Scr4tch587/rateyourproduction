import Link from "next/link";
import { apiFetch } from "@/lib/api";
import { WorkCard } from "@/components/work-card";
import type { WorkCard as WorkCardType, PaginatedResponse, LogEntry } from "@/lib/types";
import { LogEntryItem } from "@/components/log-entry";

async function getTopWorks() {
  try {
    return await apiFetch<PaginatedResponse<WorkCardType>>("/api/v1/works?per_page=10");
  } catch {
    return null;
  }
}

async function getRecentLogs() {
  try {
    return await apiFetch<PaginatedResponse<LogEntry>>("/api/v1/logs?per_page=10");
  } catch {
    return null;
  }
}

export default async function Home() {
  const [works, logs] = await Promise.all([getTopWorks(), getRecentLogs()]);

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <section className="mb-8">
        <h1 className="font-serif text-2xl font-bold text-burgundy">
          RateYourProduction
        </h1>
        <p className="text-xs text-muted-foreground mt-1">
          A catalogue of plays, musicals, and operas &mdash; logged, rated, and discovered.
        </p>
      </section>

      <div className="grid grid-cols-1 md:grid-cols-[1fr_340px] gap-10">
        <section>
          <div className="flex items-baseline justify-between mb-3 border-b border-border pb-1">
            <h2 className="font-serif text-sm font-bold">Top Rated Works</h2>
            <Link
              href="/discover"
              className="text-[10px] text-burgundy hover:text-burgundy-light transition-colors"
            >
              Browse all &rarr;
            </Link>
          </div>
          {works?.data && works.data.length > 0 ? (
            <div>
              {works.data.map((work) => (
                <WorkCard key={work.id} work={work} />
              ))}
            </div>
          ) : (
            <div className="text-xs text-muted-foreground py-6 italic">
              No works on record. The API may not be running.
            </div>
          )}
        </section>

        <aside>
          <div className="border-b border-border pb-1 mb-3">
            <h2 className="font-serif text-sm font-bold">Recent Activity</h2>
          </div>
          {logs?.data && logs.data.length > 0 ? (
            <div>
              {logs.data.map((log) => (
                <LogEntryItem key={log.id} log={log} />
              ))}
            </div>
          ) : (
            <div className="text-xs text-muted-foreground py-6 italic">
              No activity yet.
            </div>
          )}
        </aside>
      </div>
    </div>
  );
}
