package db

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

func SetUpDBTriggers(db *gorm.DB) error {

	// first execute the trigger funtion
	if db.Exec(cartTotalPriceUpdateSqlFunc).Error != nil {
		return errors.New("faild to execute cart_total_Pirce update trigger fun()")
	}

	// create trigger for calling the update funciton
	if db.Exec(cartTotalPriceTriggerExec).Error != nil {
		return errors.New("faild to create trigger for update total price on cart")
	}

	log.Printf("successfully trgger for databse are updated")
	return nil
}

var (
	// function which return total price calculation on cart when product_item added or remove delete cart
	// in here cheking  first it delete any row from cart_item then take its cart_id an find all cart_items with this id and calculate total price and update it
	// cart with this cart_id
	// any other like update or inset then take the cart_id and calculate the total price update on cart
	cartTotalPriceUpdateSqlFunc = `CREATE OR REPLACE FUNCTION update_cart_total_price() 
	RETURNS TRIGGER AS $$ 
	BEGIN 
	IF (TG_OP = 'DELETE') THEN 
		UPDATE carts c 
			SET total_price = ( 
				SELECT COALESCE ( SUM ( CASE WHEN pi.discount_price > 0 THEN pi.discount_price * ci.qty ELSE pi.price * ci.qty END), 0)::bigint 
				FROM cart_items ci INNER JOIN product_items pi ON ci.product_item_id = pi.id 
				WHERE ci.cart_id = OLD.cart_id  
			) 
		WHERE c.cart_id = OLD.cart_id; 
		RETURN NEW; 
	ELSE 
		UPDATE carts c 
			SET total_price = (
				SELECT SUM (CASE WHEN pi.discount_price > 0 THEN pi.discount_price * ci.qty ELSE pi.price * ci.qty END) 
				FROM cart_items ci INNER JOIN product_items pi ON ci.product_item_id = pi.id 
				WHERE ci.cart_id = NEW.cart_id 
			) 
			WHERE c.cart_id = NEW.cart_id;
	
	END IF; 
	RETURN NEW; 
	END; 
	$$ LANGUAGE plpgsql;`

	// for calling the trigger function above when an event of insert or update or delte happen on cart_items
	cartTotalPriceTriggerExec = `CREATE OR REPLACE TRIGGER  update_cart_total_price
	AFTER INSERT OR UPDATE OR DELETE ON cart_items
	FOR EACH ROW EXECUTE FUNCTION update_cart_total_price();`
)
