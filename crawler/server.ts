import express, { Request, Response } from 'express';
import cors from 'cors';
import { scanForTrackers } from './scan.js';


const app = express();
app.use(cors());
app.use(express.json());


app.get('/', (req: Request, res: Response) => {
  res.send('Server is running! Use POST /scan to check a URL.');
});


app.get('/healthz', (_, res) => res.status(200).send('OK'));


app.post('/scan', async (req: Request, res: Response) => {
  const { url } = req.body;
  if (!url || typeof url !== 'string') {
    return res.status(400).json({ error: 'Missing or invalid URL' });
  }

  try {
    const result = await scanForTrackers(url);
    res.json(result);
  } catch (err) {
    console.error('Scan failed', err);
    res.status(500).type('application/json').send(JSON.stringify({ error: 'Failed to scan site' }));
  }
});

app.listen(4000, '0.0.0.0', () => {
  console.log('Crawler server running on http://localhost:4000');
});
