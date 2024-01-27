import { EIP1193Provider, EIP6963ProviderDetail } from "./wallet/types";
import { registerEIP6963 } from "./wallet/provider";

const providers: EIP6963ProviderDetail[] = [];
let provider: EIP1193Provider = undefined;

declare global {
  interface Window {
    htmx: any;
  }
}

export function onPageLoad() {
  registerEIP6963(providers);
  console.log(providers);
  window.htmx.trigger("#wallets", "EIP6963", {});
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
  console.log(providers);
  return providers;
}

function findProviderByUUID(uuid: string) {
  console.log(uuid);
  for (let i = 0; i < providers.length; i++) {
    if (providers[i].info.uuid === uuid) {
      console.log(providers[i].info.uuid);
      return providers[i].provider;
    }
  }
  return undefined;
}

function findProviderByName(name: string) {
  for (let i = 0; i < providers.length; i++) {
    if (providers[i].info.name === name) {
      console.log(providers[i].info.name);
      return providers[i].provider;
    }
  }
  return undefined;
}
