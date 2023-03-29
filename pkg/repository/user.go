package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type userDatabse struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabse{DB: DB}
}

func (c *userDatabse) FindUser(ctx context.Context, user domain.User) (domain.User, error) {
	// check id,email,phone any of then match i db
	query := `SELECT * FROM users WHERE id = ? OR email = ? OR phone = ? OR user_name = ?`
	if err := c.DB.Raw(query, user.ID, user.Email, user.Phone, user.UserName).Scan(&user).Error; err != nil {
		return user, errors.New("faild to get user")
	}
	return user, nil
}

func (c *userDatabse) FindUserExceptID(ctx context.Context, user domain.User) (domain.User, error) {
	var checkUser domain.User
	query := `SELECT * FROM users WHERE id != ? AND email = ? OR id != ? AND phone = ? OR id != ? AND user_name = ?`
	if c.DB.Raw(query, user.ID, user.Email, user.ID, user.Phone, user.ID, user.UserName).Scan(&checkUser).Error != nil {
		return checkUser, errors.New("faild to check user details")
	}

	return checkUser, nil
}

func (c *userDatabse) SaveUser(ctx context.Context, user domain.User) error {

	//save the user details
	err := c.DB.Save(&user).Error

	return err
}

func (c *userDatabse) EditUser(ctx context.Context, user domain.User) error {
	fmt.Println(user)
	query := `UPDATE users SET user_name = $1, first_name = $2, last_name = $3,age = $4,email = $5, phone = $6 WHERE id = $7`
	if c.DB.Raw(query, user.UserName, user.FirstName, user.LastName, user.Age, user.Email, user.Phone, user.ID).Scan(&user).Error != nil {
		return errors.New("faild to update user")
	}
	return nil
}

// to get productItem id
func (c *userDatabse) FindProductItem(ctx context.Context, productItemID uint) (domain.ProductItem, error) {

	var productItem domain.ProductItem
	if c.DB.Raw("SELECT * FROM product_items WHERE id=?", productItemID).Scan(&productItem).Error != nil {
		return domain.ProductItem{}, errors.New("faild to get productItem from database")
	}
	return productItem, nil
}

// get the cart of user
func (c *userDatabse) FindCart(ctx context.Context, userID uint) (domain.Cart, error) {
	// then check user have cart already exist or not
	var cart domain.Cart
	if c.DB.Raw("SELECT * FROM carts WHERE user_id=?", userID).Scan(&cart).Error != nil {
		return cart, errors.New("faild to get cart of user from database")
	}

	// if user have no cart then create a new cart for user
	if cart.ID == 0 {
		querry := `INSERT INTO carts (user_id,total_price) VALUES ($1,$2) RETURNING id,user_id,total_price`
		if c.DB.Raw(querry, userID, 0).Scan(&cart).Error != nil {
			return cart, errors.New("faild to create cart for user in database")
		}
	}

	return cart, nil
}

func (c *userDatabse) UpdateCartPrice(ctx context.Context, cart domain.Cart) error {

	// update cartTotal Price
	query := ` SELECT SUM(CASE WHEN pi.discount_price > 0 THEN ci.qty * pi.discount_price ELSE ci.qty * pi.price END) AS total_price	 
	FROM cart_items ci 
	JOIN carts c ON ci.cart_id=c.id JOIN product_items pi 
	ON ci.product_item_id=pi.id AND c.id=? GROUP BY cart_id`

	var TotalPrice uint
	if c.DB.Raw(query, cart.ID).Scan(&TotalPrice).Error != nil {
		return errors.New("faild to find calculate total price of cart_items")
	}
	//update the total price on cart
	query = `UPDATE carts SET total_price = $1 WHERE id=$2`
	if c.DB.Exec(query, TotalPrice, cart.ID).Error != nil {
		return errors.New("faild to update the total price of cart")
	}

	return nil

}

// find a cartItem
func (c *userDatabse) FindCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error) {

	query := `SELECT * FROM cart_items WHERE id = ? OR cart_id=? AND product_item_id=?`
	if c.DB.Raw(query, cartItem.ID, cartItem.CartID, cartItem.ProductItemID).Scan(&cartItem).Error != nil {
		return cartItem, errors.New("faild to get cartItem of user")
	}
	return cartItem, nil
}

// add a productItem to cartitem
func (c *userDatabse) SaveCartItem(ctx context.Context, cartItems domain.CartItem) error {

	querry := `INSERT INTO cart_items (cart_id,product_item_id,qty) VALUES ($1,$2,$3)`
	if c.DB.Exec(querry, cartItems.CartID, cartItems.ProductItemID, 1).Error != nil {
		return errors.New("faild to save cart_items")
	}

	return nil
}

func (c *userDatabse) RemoveCartItem(ctx context.Context, cartItem domain.CartItem) error {

	// delete productItem from cart
	if c.DB.Raw("DELETE FROM cart_items WHERE id=?", cartItem.ID).Scan(&cartItem).Error != nil {
		return errors.New("faild to delete cart_item")
	}

	return nil
}

func (c *userDatabse) UpdateCartItem(ctx context.Context, cartItem domain.CartItem) error {

	if c.DB.Raw("UPDATE cart_items SET qty = ? WHERE id=? RETURNING qty", cartItem.Qty, cartItem.ID).Scan(&cartItem).Error != nil {
		return errors.New("faild to update the qty of cartItem")
	}
	return nil
}

// get all itmes from cart
func (c *userDatabse) GetCartItems(ctx context.Context, userId uint) (res.ResponseCart, error) {

	var response res.ResponseCart
	// get the cart of user
	cart, err := c.FindCart(ctx, userId)
	if err != nil {
		return response, err
	}

	// get the cartItem of all user with subtotal
	query := `SELECT ci.product_item_id, p.product_name, ci.qty,pi.price ,
	 pi.discount_price, CASE WHEN pi.discount_price > 0 THEN pi.discount_price * ci.qty ELSE pi.price * ci.qty END AS sub_total,  
	 pi.qty_in_stock 
	 FROM cart_items ci JOIN product_items pi ON ci.product_item_id = pi.id 
	JOIN products p ON pi.product_id=p.id AND ci.cart_id=?`

	if c.DB.Raw(query, cart.ID).Scan(&response.CartItems).Error != nil {
		return response, errors.New("faild to get cartItems from database")
	}

	response.TotalPrice = cart.TotalPrice

	return response, nil
}

func (c *userDatabse) FindAddressByID(ctx context.Context, addressID uint) (domain.Address, error) {

	var address domain.Address
	query := `SELECT * FROM addresses WHERE id=?`
	if c.DB.Raw(query, addressID).Scan(&address).Error != nil {
		return address, errors.New("faild to find address")
	}

	return address, nil
}

// find address with userId and addressess
func (c *userDatabse) FindAddressByUserID(ctx context.Context, address domain.Address, userID uint) (domain.Address, error) {

	// find the address with house,land_mark,pincode,coutry_id
	query := `SELECT * FROM addresses adrs JOIN user_addresses usr_adrs 
	ON adrs.id = usr_adrs.address_id AND user_id = ? AND house=? 
	AND land_mark=? AND pincode=? AND country_id=?`
	if c.DB.Raw(query, userID, address.House, address.LandMark, address.Pincode, address.CountryID).Scan(&address).Error != nil {
		return address, errors.New("faild to find the address")
	}
	return address, nil
}

func (c *userDatabse) FindAllAddressByUserID(ctx context.Context, userID uint) ([]res.ResAddress, error) {

	var addresses []res.ResAddress

	query := `SELECT a.id, a.house,a.name,a.phone_number,a.area,a.land_mark,a.city,a.pincode,a.country_id,c.country_name,ua.is_default
	 FROM user_addresses ua JOIN addresses a ON ua.address_id=a.id 
	 INNER JOIN countries c ON a.country_id=c.id AND ua.user_id=?`
	if c.DB.Raw(query, userID).Scan(&addresses).Error != nil {
		return addresses, errors.New("faild to get address of user")
	}

	return addresses, nil
}

func (c *userDatabse) FindCountryByID(ctx context.Context, countryID uint) (domain.Country, error) {

	var country domain.Country

	if c.DB.Raw("SELECT * FROM countries WHERE id = ?", countryID).Scan(&country).Error != nil {
		return country, errors.New("faild to find the country")
	}

	return country, nil
}

// save address
func (c *userDatabse) SaveAddress(ctx context.Context, address domain.Address) (domain.Address, error) {

	query := `INSERT INTO addresses (name,phone_number,house,area,land_mark,city,pincode,country_id) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	if c.DB.Raw(query, address.Name, address.PhoneNumber,
		address.House, address.Area, address.LandMark, address.City,
		address.Pincode, address.CountryID,
	).Scan(&address).Error != nil {
		return address, errors.New("faild to insert address on database")
	}
	return address, nil
}

// update address
func (c *userDatabse) UpdateAddress(ctx context.Context, address domain.Address) error {

	query := `UPDATE addresses SET name=$1, phone_number=$2, house=$3, area=$4, land_mark=$5, city=$6, pincode=$7,country_id=$8 WHERE id=$9`
	if c.DB.Raw(query, address.Name, address.PhoneNumber, address.House,
		address.Area, address.LandMark, address.City, address.Pincode,
		address.CountryID, address.ID).Scan(&address).Error != nil {
		return errors.New("faild to update the address for edit address")
	}
	return nil
}

func (c *userDatabse) SaveUserAddress(ctx context.Context, userAddress domain.UserAddress) (domain.UserAddress, error) {

	// if not exist then save it
	//if this address is user need to default then change old default addres to normal
	if userAddress.IsDefault {

		query := `UPDATE user_addresses SET is_default = 'f' WHERE user_id = ?`
		if c.DB.Raw(query, userAddress.UserID).Scan(&userAddress).Error != nil {
			return userAddress, errors.New("faild to remove default status of address")
		}
	}

	query := `INSERT INTO user_addresses (user_id,address_id,is_default) VALUES ($1,$2,$3) RETURNING user_id,address_id,is_default`
	if c.DB.Raw(query, userAddress.UserID, userAddress.AddressID, userAddress.IsDefault).Scan(&userAddress).Error != nil {
		return userAddress, errors.New("faild to inser userAddress on database")
	}
	return userAddress, nil
}

func (c *userDatabse) UpdateUserAddress(ctx context.Context, userAddress domain.UserAddress) error {

	// if it need to set default the change the old default
	if userAddress.IsDefault {

		query := `UPDATE user_addresses SET is_default = 'f' WHERE user_id = ?`
		if c.DB.Raw(query, userAddress.UserID).Scan(&userAddress).Error != nil {
			return errors.New("faild to remove default status of address")
		}
	}

	// update the user address
	query := `UPDATE user_addresses SET is_default = ? WHERE address_id=? AND user_id=?`
	if c.DB.Raw(query, userAddress.IsDefault, userAddress.AddressID, userAddress.UserID).Scan(&userAddress).Error != nil {
		return errors.New("faild to update user address")
	}
	return nil
}

// wish list

func (c *userDatabse) FindWishListItem(ctx context.Context, productID, userID uint) (domain.WishList, error) {

	var wishList domain.WishList
	query := `SELECT * FROM wish_lists WHERE user_id=? AND product_item_id=?`
	if c.DB.Raw(query, userID, productID).Scan(&wishList).Error != nil {
		return wishList, errors.New("faild to find wishlist item")
	}
	return wishList, nil
}

func (c *userDatabse) FindAllWishListItemsByUserID(ctx context.Context, userID uint) ([]res.ResWishList, error) {

	var wishLists []res.ResWishList

	query := `SELECT * FROM product_items pi JOIN products p ON pi.product_id=p.id JOIN wish_lists w ON w.product_item_id=pi.id AND w.user_id=?`
	if c.DB.Raw(query, userID).Scan(&wishLists).Error != nil {
		return wishLists, errors.New("faild to get wish_list items")
	}
	return wishLists, nil
}

func (c *userDatabse) SaveWishListItem(ctx context.Context, wishList domain.WishList) error {

	query := `INSERT INTO wish_lists (user_id,product_item_id) VALUES ($1,$2) RETURNING *`

	if c.DB.Raw(query, wishList.UserID, wishList.ProductItemID).Scan(&wishList).Error != nil {
		return errors.New("faild to insert new wishlist on database")
	}
	return nil
}

func (c *userDatabse) RemoveWishListItem(ctx context.Context, wishList domain.WishList) error {

	query := `DELETE FROM wish_lists WHERE id=?`
	if c.DB.Raw(query, wishList.ID).Scan(&wishList).Error != nil {
		return errors.New("faild to delete productItem from database")
	}
	return nil
}

// checkout page
func (c *userDatabse) CheckOutCart(ctx context.Context, userId uint) (res.ResCheckOut, error) {

	var resCheckOut res.ResCheckOut
	// get all cartItems of user which are not out of stock
	query := `SELECT ci.product_item_id, p.product_name,pi.price,pi.discount_price, pi.qty_in_stock, ci.qty, 
	CASE WHEN pi.discount_price > 0 THEN (ci.qty * pi.discount_price) ELSE (ci.qty * pi.price) END AS sub_total  
	FROM cart_items ci JOIN carts c ON ci.cart_id = c.id JOIN product_items pi ON ci.product_item_id = pi.id 
	JOIN products p ON pi.product_id = p.id AND c.user_id = ?`

	if c.DB.Raw(query, userId).Scan(&resCheckOut.ProductItems).Error != nil {
		return resCheckOut, errors.New("faild to get cartItems for checkout")
	}

	// get user addresses
	adresses, err := c.FindAllAddressByUserID(ctx, userId)
	if err != nil {
		return resCheckOut, errors.New("faild to get user addrss for checkout")
	}
	resCheckOut.Addresses = adresses

	// find total price
	query = `SELECT total_price FROM carts WHERE user_id = ?`
	if c.DB.Raw(query, userId).Scan(&resCheckOut.TotalPrice).Error != nil {
		return resCheckOut, errors.New("faild to ge total price for user cart")
	}

	return resCheckOut, nil
}
