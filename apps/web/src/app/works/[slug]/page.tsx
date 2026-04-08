import { notFound } from "next/navigation";
import { LogForm } from "@/components/log-form";
import { Rating } from "@/components/rating";
import { ProductionTable } from "@/components/production-table";
import { Separator } from "@/components/ui/separator";
import { apiFetch } from "@/lib/api";
import type { WorkDetail, PaginatedResponse, LogEntry } from "@/lib/types";
import { LogEntryItem } from "@/components/log-entry";
import { SubmissionForm } from "@/components/submission-form";

async function getWork(slug: string) {
  try {
    return await apiFetch<WorkDetail>(`/api/v1/works/${slug}`);
  } catch {
    return null;
  }
}

async function getWorkLogs(workId: string) {
  try {
    return await apiFetch<PaginatedResponse<LogEntry>>(`/api/v1/logs?work_id=${workId}&per_page=10`);
  } catch {
    return null;
  }
}

interface WorkPageProps {
  params: Promise<{ slug: string }>;
}

export default async function WorkPage({ params }: WorkPageProps) {
  const { slug } = await params;
  const work = await getWork(slug);

  if (!work) {
    notFound();
  }

  const logs = await getWorkLogs(work.id);

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="grid grid-cols-1 md:grid-cols-[280px_1fr] gap-10">
        {/* Left column — margin notes style */}
        <aside className="space-y-4">
          <div>
            <h1 className="font-serif text-xl font-bold leading-tight">
              {work.title}
            </h1>
            <div className="text-xs italic text-muted-foreground mt-0.5">
              {work.type}
            </div>
          </div>

          {work.creators.length > 0 && (
            <div>
              {work.creators.map((c) => (
                <div key={c.person_id} className="text-xs">
                  <span className="font-medium">{c.name}</span>
                  <span className="text-muted-foreground italic"> {c.role_type}</span>
                </div>
              ))}
            </div>
          )}

          <div>
            <Rating rating={work.average_rating} count={work.rating_count} />
          </div>

          {work.genres.length > 0 && (
            <div className="flex flex-wrap gap-1.5">
              {work.genres.map((g) => (
                <span
                  key={g.id}
                  className="text-[10px] border border-border px-1.5 py-0.5 text-muted-foreground"
                >
                  {g.name}
                </span>
              ))}
            </div>
          )}

          {work.premiere_year && (
            <div className="text-xs text-muted-foreground">
              Premiered {work.premiere_year}
            </div>
          )}

          {work.description && (
            <>
              <Separator />
              <p className="text-xs text-muted-foreground leading-relaxed italic">
                {work.description}
              </p>
            </>
          )}
        </aside>

        {/* Right column — main document */}
        <div>
          <div className="mb-6 grid gap-4 lg:grid-cols-2">
            <LogForm workId={work.id} productions={work.productions || []} />
            <SubmissionForm workId={work.id} />
          </div>

          <div className="border-b border-border pb-1 mb-3">
            <h2 className="font-serif text-sm font-bold">Productions</h2>
          </div>
          <ProductionTable productions={work.productions || []} />

          <div className="mt-8">
            <div className="border-b border-border pb-1 mb-3">
              <h2 className="font-serif text-sm font-bold">Recent Logs</h2>
            </div>
            {logs?.data && logs.data.length > 0 ? (
              <div>
                {logs.data.map((log) => (
                  <LogEntryItem key={log.id} log={log} />
                ))}
              </div>
            ) : (
              <div className="text-xs text-muted-foreground py-4 italic">
                No logs yet. Be the first to record this work.
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
