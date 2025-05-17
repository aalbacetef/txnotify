export type Transaction = {
  hash: string;
  from: string;
  to: string;
  value: bigint;
}

export type Endpoint = {
  url: string;
}

export type ChainlistDataEntry = {
  name: string;
  chain: string;
  rpc: Endpoint[];
}

export async function getEndpoints(): Promise<Endpoint[]> {
  const response = await fetch(
    'https://chainlist.org/rpcs.json',
    { method: 'GET', mode: 'cors' },
  );
  const data = await response.json() as ChainlistDataEntry[];

  const mainnet = data.find(
    elem => elem.name === "Ethereum Mainnet" && elem.chain === "ETH"
  );

  if (typeof mainnet === 'undefined') {
    console.error("could not find RPC endpoints");
    return [];
  }


  return mainnet.rpc.filter(elem => elem.url.startsWith("https://"));
}
