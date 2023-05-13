package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
	"gorm.io/gorm"
)

type userDatabse struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabse{DB: DB}
}
func (c *userDatabse) FindUserByUserID(ctx context.Context, userID uint) (user domain.User, err error) {

	query := `SELECT * FROM users WHERE id = $1`
	err = c.DB.Raw(query, userID).Scan(&user).Error

	return user, err
}

func (c *userDatabse) FindUserByEmail(ctx context.Context, email string) (user domain.User, err error) {

	query := `SELECT * FROM users WHERE email = $1`
	err = c.DB.Raw(query, email).Scan(&user).Error

	return user, err
}

func (c *userDatabse) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (user domain.User, err error) {

	query := `SELECT * FROM users WHERE phone = $1`
	err = c.DB.Raw(query, phoneNumber).Scan(&user).Error

	return user, err
}
func (c *userDatabse) FindUserByUserName(ctx context.Context, userName string) (user domain.User, err error) {

	query := `SELECT * FROM users WHERE user_name = $1`
	err = c.DB.Raw(query, userName).Scan(&user).Error

	return user, err
}

func (c *userDatabse) CheckOtherUserWithDetails(ctx context.Context, user domain.User) (domain.User, error) {
	var checkUser domain.User
	query := `SELECT * FROM users WHERE id != ? AND email = ? OR id != ? AND phone = ? OR id != ? AND user_name = ?`
	if c.DB.Raw(query, user.ID, user.Email, user.ID, user.Phone, user.ID, user.UserName).Scan(&checkUser).Error != nil {
		return checkUser, errors.New("faild to check user details")
	}

	return checkUser, nil
}

func (c *userDatabse) SaveUser(ctx context.Context, user domain.User) (userID uint, err error) {

	//save the user details
	query := `INSERT INTO users (user_name, first_name, last_name, age, email, phone, password,created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	createdAt := time.Now()
	err = c.DB.Raw(query, user.UserName, user.FirstName, user.LastName,
		user.Age, user.Email, user.Phone, user.Password, createdAt).Scan(&user).Error

	return userID, err
}

func (c *userDatabse) UpdateUser(ctx context.Context, user domain.User) (err error) {

	updatedAt := time.Now()
	// check password need to update or not
	if user.Password != "" {
		query := `UPDATE users SET user_name = $1, first_name = $2, last_name = $3,age = $4, 
		email = $5, phone = $6, password = $7, updated_at = $8 WHERE id = $9`
		err = c.DB.Exec(query, user.UserName, user.FirstName, user.LastName, user.Age, user.Email,
			user.Phone, user.Password, updatedAt, user.ID).Error
	} else {
		query := `UPDATE users SET user_name = $1, first_name = $2, last_name = $3,age = $4, 
		email = $5, phone = $6, updated_at = $7 WHERE id = $8`
		err = c.DB.Exec(query, user.UserName, user.FirstName, user.LastName, user.Age, user.Email,
			user.Phone, updatedAt, user.ID).Error
	}

	if err != nil {
		return fmt.Errorf("faild to update user detail of user with user_id %d", user.ID)
	}
	return nil
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
func (c *userDatabse) SaveAddress(ctx context.Context, address domain.Address) (addressID uint, err error) {

	query := `INSERT INTO addresses (name,phone_number,house,area,land_mark,city,pincode,country_id,created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	createdAt := time.Now()

	if c.DB.Raw(query, address.Name, address.PhoneNumber,
		address.House, address.Area, address.LandMark, address.City,
		address.Pincode, address.CountryID, createdAt,
	).Scan(&address).Error != nil {
		return addressID, errors.New("faild to insert address on database")
	}
	return address.ID, nil
}

// update address
func (c *userDatabse) UpdateAddress(ctx context.Context, address domain.Address) error {

	query := `UPDATE addresses SET name=$1, phone_number=$2, house=$3, area=$4, land_mark=$5, 
	city=$6, pincode=$7,country_id=$8, updated_at = $9 WHERE id=$10`

	updatedAt := time.Now()
	if c.DB.Raw(query, address.Name, address.PhoneNumber, address.House,
		address.Area, address.LandMark, address.City, address.Pincode,
		address.CountryID, updatedAt, address.ID).Scan(&address).Error != nil {
		return errors.New("faild to update the address for edit address")
	}
	return nil
}

func (c *userDatabse) SaveUserAddress(ctx context.Context, userAddress domain.UserAddress) error {

	// first check user's first address is this or not
	var userID uint
	query := `SELECT address_id FROM user_addresses WHERE user_id = $1`
	err := c.DB.Raw(query, userAddress.UserID).Scan(&userID).Error
	if err != nil {
		return fmt.Errorf("faild to check user have already address exit or not with user_id %v", userAddress.UserID)
	}

	// if the given address is need to set default  then remove all other from default
	if userID == 0 { // it means user have no other addresses
		userAddress.IsDefault = true
	} else if userAddress.IsDefault {
		query := `UPDATE user_addresses SET is_default = 'f' WHERE user_id = ?`
		if c.DB.Raw(query, userAddress.UserID).Scan(&userAddress).Error != nil {
			return errors.New("faild to remove default status of address")
		}
	}

	query = `INSERT INTO user_addresses (user_id,address_id,is_default) VALUES ($1, $2, $3)`
	err = c.DB.Exec(query, userAddress.UserID, userAddress.AddressID, userAddress.IsDefault).Error
	if err != nil {
		return errors.New("faild to inser userAddress on database")
	}
	return nil
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
