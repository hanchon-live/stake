import { EIP1193Provider, EIP6963ProviderDetail } from "./wallet/types";
import { registerEIP6963 } from "./wallet/provider";

const providers: EIP6963ProviderDetail[] = [];
let provider: EIP1193Provider = undefined;

declare global {
  interface Window {
    htmx: any;
  }
}

export async function getAddress() {
  let accounts: string[] = [];
  const p = getProvider();
  if (p != undefined) {
    accounts = (await p.request({
      method: "eth_requestAccounts",
    })) as string[];
  }

  window.htmx.ajax("POST", "/currentwallet", {
    target: "#currentwallet",
    swap: "innerHTML",
    values: { accounts: accounts },
  });

  return accounts;
}

export function onPageLoad() {
  registerEIP6963(providers);
  // TODO: move this to one second after page loads
  const walletNames = [];
  for (let i = 0; i < providers.length; ++i) {
    walletNames.push(providers[i].info.name);
  }
  window.htmx.ajax("POST", "/wallets", {
    target: "#walletslist",
    swap: "innerHTML",
    values: { providers: walletNames },
  });
}

export function setProviderByUUID(uuid: string) {
  provider = findProviderByUUID(uuid);
}

export function setProviderByName(name: string) {
  provider = findProviderByName(name);
}

export function getProvider() {
  return provider;
}

export function getProviders() {
  return providers;
}

function findProviderByUUID(uuid: string) {
  for (let i = 0; i < providers.length; i++) {
    if (providers[i].info.uuid === uuid) {
      return providers[i].provider;
    }
  }
  return undefined;
}

function findProviderByName(name: string) {
  for (let i = 0; i < providers.length; i++) {
    if (providers[i].info.name === name) {
      return providers[i].provider;
    }
  }
  return undefined;
}
