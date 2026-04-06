import { Page_marketplace_product } from 'app/features/marketplace/Page_marketplace_product'

export default async function MarketplaceProductPage({
  params,
}: {
  params: Promise<{ slug: string }>
}) {
  const { slug } = await params
  return <Page_marketplace_product slug={slug} />
}
