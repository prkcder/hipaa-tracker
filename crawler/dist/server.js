// import express from 'express';
// import { runScan } from './scan';
// const app = express();
// app.use(express.json());
// app.post('/scan', async (req, res) => {
//     const { url } = req.body;
//     if (!url || typeof url !== 'string') {
//         return res.status(400).json({ error: 'Invalid or missing URL' });
//     }
//     try {
//         const results = await runScan(url);
//         res.json(results);
//     } catch (err) {
//         console.error('Scan error:', err);
//         res.status(500).json({ error: 'Scan failed' });
//     }
// });
// const PORT = 4000;
// app.listen(PORT, () => console.log(`Crawler server running on port ${PORT}`));
import express from 'express';
const app = express();
//# sourceMappingURL=server.js.map