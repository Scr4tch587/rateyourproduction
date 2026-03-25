interface WorkPageProps {
  params: Promise<{ slug: string }>;
}

export default async function WorkPage({ params }: WorkPageProps) {
  const { slug } = await params;

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="grid grid-cols-1 md:grid-cols-[1fr_360px] gap-8">
        <div>
          <h1 className="text-lg font-bold mb-1">{slug}</h1>
          <div className="text-xs text-muted-foreground mb-4">
            Work details will load here.
          </div>
          <div className="space-y-2">
            <div className="text-xs">
              <span className="text-muted-foreground">Type:</span> —
            </div>
            <div className="text-xs">
              <span className="text-muted-foreground">Creators:</span> —
            </div>
            <div className="text-xs">
              <span className="text-muted-foreground">Genres:</span> —
            </div>
            <div className="text-xs">
              <span className="text-muted-foreground">Rating:</span> —
            </div>
          </div>
        </div>
        <aside>
          <h2 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-3">
            Productions
          </h2>
          <div className="text-xs text-muted-foreground">
            No productions found.
          </div>
        </aside>
      </div>
    </div>
  );
}
