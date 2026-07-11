-- Retire legacy donation / purchase / subscription custom menu entries.
--
-- The CAU deployment now uses the LDXP card flow plus in-site redeem.
-- Keep the backend custom-menu feature available, but remove old visible
-- sidebar entries that can send users into inactive donation/payment screens.

DO $$
DECLARE
    v_raw      text;
    v_items    jsonb;
    v_filtered jsonb;
BEGIN
    SELECT value INTO v_raw
      FROM settings
     WHERE key = 'custom_menu_items';

    IF COALESCE(v_raw, '') = '' OR v_raw = 'null' THEN
        RETURN;
    END IF;

    v_items := v_raw::jsonb;

    SELECT COALESCE(jsonb_agg(elem ORDER BY ord), '[]'::jsonb)
      INTO v_filtered
      FROM jsonb_array_elements(v_items) WITH ORDINALITY AS t(elem, ord)
     WHERE NOT (
        lower(COALESCE(elem ->> 'id', '')) IN (
            'migrated_purchase_subscription',
            'donation',
            'donate',
            'sponsor',
            'purchase_subscription'
        )
        OR COALESCE(elem ->> 'label', '') ~* '(捐赠|赞助|打赏|收款|捐钱|充值[[:space:]]*/[[:space:]]*订阅|订阅|donat|donate|donation|sponsor|tip|purchase|subscription)'
        OR COALESCE(elem ->> 'url', '') ~* '(/purchase|donat|donate|donation|sponsor|tip)'
    );

    IF v_filtered IS DISTINCT FROM v_items THEN
        UPDATE settings
           SET value = v_filtered::text
         WHERE key = 'custom_menu_items';
    END IF;

    UPDATE settings SET value = 'false' WHERE key = 'purchase_subscription_enabled';
    UPDATE settings SET value = '' WHERE key = 'purchase_subscription_url';
END $$;
