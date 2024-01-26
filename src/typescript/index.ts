import { EIP1193Provider, EIP6963ProviderDetail } from "./wallet/types";
import { registerEIP6963 } from "./wallet/provider";

const providers: EIP6963ProviderDetail[] = [];
let provider: EIP1193Provider = undefined;

export function onPageLoad() {
  registerEIP6963(providers);
  console.log(providers);
}

export function setProvider(uuid: string) {
  provider = findProvider(uuid);
}

export function getProvider() {
  return provider;
}

export function getProviders() {
  console.log(providers);
  return providers;
}

function findProvider(uuid: string) {
  console.log(uuid);
  for (let i = 0; i < providers.length; i++) {
    if (providers[i].info.uuid === uuid) {
      console.log(providers[i].info.uuid);
      return providers[i].provider;
    }
  }
  return undefined;
}
