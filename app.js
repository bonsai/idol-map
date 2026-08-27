// idol-map app: pick region + decade -> YouTube search for that region's
// local idols, plus the decade's idol-folklore context from idol-quiz.

const { regions, decades, folklore, localIdols } = IDOL_MAP_DATA;

let selectedRegion = null;
let selectedDecade = null;

const mapEl = document.getElementById("map");
const timelineEl = document.getElementById("timeline");
const exhibitEl = document.getElementById("exhibit");
const exhibitHint = document.getElementById("exhibitHint");
const nowRegion = document.getElementById("nowRegion");
const nowDecade = document.getElementById("nowDecade");
const playBtn = document.getElementById("playBtn");

function renderMap() {
  mapEl.innerHTML = "";
  regions.forEach((r) => {
    const el = document.createElement("div");
    el.className = "region" + (selectedRegion === r.id ? " active" : "");
    el.textContent = r.name;
    el.addEventListener("click", () => {
      selectedRegion = r.id;
      renderMap();
      update();
    });
    mapEl.appendChild(el);
  });
}

function renderTimeline() {
  timelineEl.innerHTML = "";
  decades.forEach((d) => {
    const el = document.createElement("div");
    el.className = "decade" + (selectedDecade === d ? " active" : "");
    el.textContent = d + "s";
    el.addEventListener("click", () => {
      selectedDecade = d;
      renderTimeline();
      update();
    });
    timelineEl.appendChild(el);
  });
}

function regionName(id) {
  const r = regions.find((x) => x.id === id);
  return r ? r.name : id;
}

function youtubeSearch(regionId, decade) {
  const q = `御当地アイドル ${regionName(regionId)} ${decade}年代`;
  return "https://www.youtube.com/results?search_query=" + encodeURIComponent(q);
}

function update() {
  const ready = selectedRegion && selectedDecade;
  nowRegion.textContent = selectedRegion ? regionName(selectedRegion) : "地方を選択";
  nowDecade.textContent = selectedDecade ? selectedDecade + "年代" : "年代を選択";
  playBtn.disabled = !ready;
  if (ready) {
    playBtn.textContent = `▶ ${regionName(selectedRegion)} の ${selectedDecade}年代 御当地アイドルを再生`;
  }

  exhibitEl.innerHTML = "";
  exhibitHint.textContent = "";

  if (!selectedDecade) {
    exhibitEl.innerHTML = '<div class="empty">年代を選ぶと、その時代のアイドル民俗学が表示されます。</div>';
    return;
  }

  const pins = (localIdols || []).filter(
    (x) => x.region === selectedRegion && x.decade === selectedDecade
  );
  if (selectedRegion && pins.length) {
    exhibitHint.textContent = "・この地域のピン";
    const wrap = document.createElement("div");
    wrap.className = "region-pins";
    pins.forEach((p) => {
      const a = document.createElement("a");
      a.className = "pin";
      a.href = "https://www.youtube.com/results?search_query=" + encodeURIComponent(p.name);
      a.target = "_blank";
      a.rel = "noopener";
      a.textContent = "📍 " + p.name;
      wrap.appendChild(a);
    });
    exhibitEl.appendChild(wrap);
  }

  const items = folklore[selectedDecade] || [];
  if (!items.length) {
    exhibitEl.insertAdjacentHTML("beforeend", '<div class="empty">データなし</div>');
    return;
  }
  items.forEach((it) => {
    const card = document.createElement("div");
    card.className = "card";
    card.innerHTML = `<div class="q">${it.q}</div><div class="a">${it.a}</div>`;
    exhibitEl.appendChild(card);
  });
}

playBtn.addEventListener("click", () => {
  if (!selectedRegion || !selectedDecade) return;
  window.open(youtubeSearch(selectedRegion, selectedDecade), "_blank", "noopener");
});

renderMap();
renderTimeline();
update();
