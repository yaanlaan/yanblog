
const articles = [
    { CreatedAt: "2026-02-01T10:00:00+08:00" },
    { CreatedAt: "2026-02-01T12:00:00+08:00" },
    { CreatedAt: "2026-02-01T14:00:00+08:00" },
    { CreatedAt: "2026-02-01T16:00:00+08:00" },
    { CreatedAt: "2026-02-01T18:00:00+08:00" },
];

const toDateKey = (date) => {
  const y = date.getFullYear()
  const m = (date.getMonth() + 1).toString().padStart(2, '0')
  const d = date.getDate().toString().padStart(2, '0')
  return `${y}-${m}-${d}`
}

// 模拟 contributionData
const contributionData = new Map()
articles.forEach(art => {
    const dStr = art.CreatedAt
    if (dStr) {
      const date = new Date(dStr)
      if (!isNaN(date.getTime())) {
         const key = toDateKey(date)
         contributionData.set(key, (contributionData.get(key) || 0) + 1)
      }
    }
})

console.log('Contribution Data:', Object.fromEntries(contributionData));

// 模拟 calendarWeeks
const weeks = []
// 模拟当前时间是 2026-02-01 (注意：Node 环境下 Date() 是当前的真实时间，
// 但为了复现用户的 2026-02-01，我们需要 Mock 今天的日期，或者假设系统时间已经被修改)
// 我们可以手动指定 start 的基准点为 2026-02-01
const now = new Date("2026-02-01T20:00:00+08:00"); 

const start = new Date(now)
start.setFullYear(start.getFullYear() - 1)

// 调整到该周的周日
while (start.getDay() !== 0) {
    start.setDate(start.getDate() - 1)
}

const current = new Date(start)
console.log('Grid Start Date:', toDateKey(current));

let found = false;
for (let w = 0; w < 53; w++) {
    for (let i = 0; i < 7; i++) {
        const dateStr = toDateKey(current)
        const count = contributionData.get(dateStr) || 0
        
        if (dateStr === "2026-02-01") {
            console.log(`FOUND 2026-02-01 at week ${w}, day ${i}, count: ${count}`);
            found = true;
        }

        current.setDate(current.getDate() + 1)
    }
}

if (!found) {
    console.log("2026-02-01 NOT FOUND in grid!");
    console.log("Grid End Date:", toDateKey(current));
}
