package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post( post_id, title, content, author_id, community_id ) 
				values (?,?,?,?,?) `
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostById 根据id查询单个帖子数据详情
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select 
			post_id, title, content, author_id, community_id, create_time 
			from post
			where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 查询帖子列表函数（帖子从新到旧的顺序返回）
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select 
			post_id, title, content, author_id, community_id, create_time 
			from post
			ORDER BY create_time
			DESC
			limit ?,? `
	posts = make([]*models.Post, 0, 2) // make创建切片，长度为0，容量为5
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// 根据给定的id列表查询帖子数据
func GetPostByIDs(ids []string) (postList []*models.Post, err error) {
	// order by FIND_IN_SET(post_id, ?) 确保 MySQL 按这个 ids 的顺序来排
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
			   from post
			   where post_id in (?)
			   order by FIND_IN_SET(post_id, ?)`
	// strings.Join 将切片里的每个元素按顺序连接起来，并用分隔符隔开，返回一个新的字符串
	// (?) 中的参数量是不固定的，可能是1个或者10个，需要 sqlx.In 来展开切片
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	// db.Rebind 会根据你当前数据库驱动的需求，把 query 里面的占位符 ? 转换成该驱动支持的格式。
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
