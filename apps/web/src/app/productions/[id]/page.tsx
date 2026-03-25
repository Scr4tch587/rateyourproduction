interface ProductionPageProps {
  params: Promise<{ id: string }>;
}

export default async function ProductionPage({ params }: ProductionPageProps) {
  const { id } = await params;

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="text-xs text-muted-foreground mb-2">
        ← Back to work
      </div>
      <h1 className="text-lg font-bold mb-1">Production #{id}</h1>
      <div className="grid grid-cols-1 md:grid-cols-[1fr_360px] gap-8 mt-4">
        <div className="space-y-2 text-xs">
          <div>
            <span className="text-muted-foreground">Company:</span> —
          </div>
          <div>
            <span className="text-muted-foreground">Venue:</span> —
          </div>
          <div>
            <span className="text-muted-foreground">Dates:</span> —
          </div>
          <div>
            <span className="text-muted-foreground">Rating:</span> —
          </div>
        </div>
        <aside>
          <h2 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-3">
            Recent Logs
          </h2>
          <div className="text-xs text-muted-foreground">No logs yet.</div>
        </aside>
      </div>
    </div>
  );
}
