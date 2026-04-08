import Link from "next/link";
import { notFound } from "next/navigation";
import { LogForm } from "@/components/log-form";
import { Rating } from "@/components/rating";
import { apiFetch } from "@/lib/api";
import type { ProductionDetail, PaginatedResponse, LogEntry, Production } from "@/lib/types";
import { LogEntryItem } from "@/components/log-entry";

async function getProduction(id: string) {
  try {
    return await apiFetch<ProductionDetail>(`/api/v1/productions/${id}`);
  } catch {
    return null;
  }
}

async function getProductionLogs(productionId: string) {
  try {
    return await apiFetch<PaginatedResponse<LogEntry>>(`/api/v1/logs?production_id=${productionId}&per_page=10`);
  } catch {
    return null;
  }
}

interface ProductionPageProps {
  params: Promise<{ id: string }>;
}

export default async function ProductionPage({ params }: ProductionPageProps) {
  const { id } = await params;
  const production = await getProduction(id);

  if (!production) {
    notFound();
  }

  const logs = await getProductionLogs(production.id);

  const directors = production.credits?.filter((c) => c.role_type === "director") || [];
  const actors = production.credits?.filter((c) => c.role_type === "actor") || [];
  const otherCredits = production.credits?.filter((c) => c.role_type !== "director" && c.role_type !== "actor") || [];

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <Link
        href={`/works/${production.work_slug}`}
        className="text-xs text-burgundy hover:text-burgundy-light transition-colors mb-4 block"
      >
        &larr; {production.work_title}
      </Link>

      <div className="grid grid-cols-1 md:grid-cols-[280px_1fr] gap-10">
        {/* Left column — production details */}
        <aside className="space-y-4">
          <h1 className="font-serif text-xl font-bold leading-tight">
            {production.production_label || production.company_name || production.city || "Production"}
          </h1>

          <Rating rating={production.average_rating} count={production.rating_count} />

          <div className="space-y-1.5 text-xs">
            {production.company_name && (
              <div>
                <span className="text-muted-foreground">Company</span>
                <div className="font-medium">{production.company_name}</div>
              </div>
            )}
            {production.venue_name && (
              <div>
                <span className="text-muted-foreground">Venue</span>
                <div className="font-medium">{production.venue_name}</div>
              </div>
            )}
            {production.city && (
              <div>
                <span className="text-muted-foreground">Location</span>
                <div className="font-medium">
                  {production.city}{production.country && `, ${production.country}`}
                </div>
              </div>
            )}
            {production.start_date && (
              <div>
                <span className="text-muted-foreground">Dates</span>
                <div className="font-medium">
                  {production.start_date}{production.end_date && ` \u2014 ${production.end_date}`}
                </div>
              </div>
            )}
          </div>
        </aside>

        {/* Right column — credits + logs */}
        <div>
          <div className="mb-6">
            <LogForm
              workId={production.work_id}
              productions={[production as Production]}
              defaultProductionId={production.id}
            />
          </div>

          {directors.length > 0 && (
            <section className="mb-6">
              <div className="border-b border-border pb-1 mb-2">
                <h2 className="font-serif text-sm font-bold">
                  Director{directors.length > 1 ? "s" : ""}
                </h2>
              </div>
              <div className="text-xs space-y-0.5">
                {directors.map((c) => (
                  <div key={c.person_id} className="font-medium">{c.name}</div>
                ))}
              </div>
            </section>
          )}

          {actors.length > 0 && (
            <section className="mb-6">
              <div className="border-b border-border pb-1 mb-2">
                <h2 className="font-serif text-sm font-bold">Cast</h2>
              </div>
              <div className="text-xs space-y-0.5">
                {actors.map((c) => (
                  <div key={c.person_id} className="font-medium">{c.name}</div>
                ))}
              </div>
            </section>
          )}

          {otherCredits.length > 0 && (
            <section className="mb-6">
              <div className="border-b border-border pb-1 mb-2">
                <h2 className="font-serif text-sm font-bold">Creative Team</h2>
              </div>
              <div className="text-xs space-y-0.5">
                {otherCredits.map((c) => (
                  <div key={c.person_id + c.role_type}>
                    <span className="font-medium">{c.name}</span>{" "}
                    <span className="text-muted-foreground italic">{c.role_type}</span>
                  </div>
                ))}
              </div>
            </section>
          )}

          <section>
            <div className="border-b border-border pb-1 mb-3">
              <h2 className="font-serif text-sm font-bold">Logs</h2>
            </div>
            {logs?.data && logs.data.length > 0 ? (
              <div>
                {logs.data.map((log) => (
                  <LogEntryItem key={log.id} log={log} />
                ))}
              </div>
            ) : (
              <div className="text-xs text-muted-foreground py-4 italic">
                No logs yet.
              </div>
            )}
          </section>
        </div>
      </div>
    </div>
  );
}
