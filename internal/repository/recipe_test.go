package repository_test

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"vutung2311-golang-test/internal/repository"
	"vutung2311-golang-test/pkg/httpclient"
	"vutung2311-golang-test/pkg/worker"
)

type response struct {
	statusCode   int
	responseData string
}

func TestRecipeRepository_FindByIDs(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(new(logrus.JSONFormatter))

	recipeResponseMap := map[string]response{
		"1": {
			statusCode:   http.StatusOK,
			responseData: `{"id":"1","name":"Parmesan-Crusted Pork Tenderloin","headline":"with Potato Wedges and Apple Walnut Salad","description":"Parm\u2019s the charm with this next-level pork recipe. The cheese is mixed with panko breadcrumbs to create a crust that coats the tenderloin like a glorious golden-brown crown. That way, you get meltiness, juiciness, and crunch in every bite. But this recipe isn\u2019t just about the meat: there\u2019s also roasted rosemary potatoes and a crisp apple walnut salad to round things out.","difficulty":1,"prepTime":"PT30M","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/0,0\/image\/parmesan-crusted-pork-tenderloin-66608000.jpg","ingredients":[{"name":"Rosemary","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/55661a71f8b25e391e8b456a.png"},{"name":"Yukon Gold Potatoes","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/554a3abff8b25e1d268b456d.png"},{"name":"Parmesan Cheese","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5550e133fd2cb9a7168b456b.png"},{"name":"Garlic Powder","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/55ef01fbf8b25eba7e8b4567.png"},{"name":"Panko Breadcrumbs","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/ingredients\/554a39a04dab71626c8b456b-3f519176.png"},{"name":"Pork Tenderloin","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5567235df8b25e472f8b4567.png"},{"name":"Sour Cream","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5550e1064dab71893e8b4569.png"},{"name":"Lemon","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/554a302ffd2cb9324b8b4569.png"},{"name":"Apple","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/554a3a7cf8b25ed7288b456b.png"},{"name":"Spring Mix Lettuce","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5566e35f4dab71ea078b4567.png"},{"name":"Dried Cranberries","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5626b14af8b25e0b1f8b4567.png"},{"name":"Walnuts","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5550df2afd2cb9dd178b4569.png"},{"name":"Vegetable Oil","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/ingredients\/5566d4f94dab715a078b4568-7c93a003.png"},{"name":"Olive Oil","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5566cdf2f8b25e0d298b4568.png"},{"name":"Salt","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5566ceb7fd2cb95f7f8b4567.png"},{"name":"Pepper","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5566dc00f8b25e5b298b4568.png"}]}`,
		},
		"2": {
			statusCode:   http.StatusOK,
			responseData: `{"id":"2","name":"Melty Monterey Jack Burgers","headline":"with Red Onion Jam and Zucchini Fries","description":"There are a lot of burger recipes out there, we know. But we like to think that this is the burger that tops all burgers, thanks to oozy, melty Monterey Jack cheese and jammy balsamic onions. Oh, and we should mention that this patty comes with breaded zucchini on the side, so you can get your fix of crispy, crunchy finger foods and a dose of veg all in one.","difficulty":1,"prepTime":"PT35M","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/0,0\/image\/5a958c3830006c344a20aba2-3e75241e.jpg","ingredients":[{"name":"Garlic","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/554a363df8b25e1d268b456b.png"},{"name":"Red Onion","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/554a3947fd2cb9cf488b456b.png"},{"name":"Zucchini","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5553981df8b25e5d0c8b456a.png"},{"name":"Mayonnaise","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/554a2f494dab71636c8b4569.png"},{"name":"Balsamic Vinegar","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/554a3879fd2cb9ba4f8b456a.png"},{"name":"Panko Breadcrumbs","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/ingredients\/554a39a04dab71626c8b456b-3f519176.png"},{"name":"Dried Oregano","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5566de3b4dab715b078b456a.png"},{"name":"Ground Beef","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/554a35754dab71d76f8b4568.png"},{"name":"Monterey Jack Cheese","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/55c109e74dab7112098b4567.png"},{"name":"Potato Buns","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/561d2ffff8b25eb4118b4567.png"},{"name":"Ketchup","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/556715a1f8b25ea22e8b4569.png"},{"name":"Vegetable Oil","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/ingredients\/5566d4f94dab715a078b4568-7c93a003.png"},{"name":"Sugar","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5566e2014dab71e3078b4568.png"},{"name":"Salt","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5566ceb7fd2cb95f7f8b4567.png"},{"name":"Pepper","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5566dc00f8b25e5b298b4568.png"}]}`,
		},
		"3": {
			statusCode:   http.StatusOK,
			responseData: `{"id":"3","name":"Tex-Mex Tilapia","headline":"with Cilantro Lime Couscous and Green Beans","description":"Let\u2019s take tilapia to the next level and turn it into a Tex-Mex-style triumph on your plate. The firm-fleshed fish is given a dusting of panko breadcrumbs and our Southwest spice blend, which ensures that it has satisfying crunch and zesty flavor in every bite. The couscous, green beans, and lime crema on the side come together in a flash, meaning you\u2019ll have this masterpiece of a meal on the table in a matter of minutes.","difficulty":1,"prepTime":"PT20M","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/0,0\/image\/5a958c0d30006c33ca2850f2-c352c2d5.jpg","ingredients":[{"name":"Vegetable Stock Concentrate","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/554a3a9e4dab71716c8b456b.png"},{"name":"Cilantro","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/554a3c4d4dab71716c8b456c.png"},{"name":"Lime","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/ingredient\/554a3c9efd2cb9ba4f8b456c-f32287bd.png"},{"name":"Panko Breadcrumbs","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/ingredients\/554a39a04dab71626c8b456b-3f519176.png"},{"name":"Southwest Spice Blend","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/57484a434dab718c228b4567.png"},{"name":"Couscous","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5550dab4f8b25e56468b456c.png"},{"name":"Tilapia","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5566dd92f8b25eb4298b4567.png"},{"name":"Sour Cream","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5550e1064dab71893e8b4569.png"},{"name":"Chipotle Powder","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/554a3c424dab71d76f8b456a.png"},{"name":"Green Beans","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5553931dfd2cb9db798b4569.png"},{"name":"Vegetable Oil","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/ingredients\/5566d4f94dab715a078b4568-7c93a003.png"},{"name":"Salt","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5566ceb7fd2cb95f7f8b4567.png"},{"name":"Pepper","imageLink":"https:\/\/d3hvwccx09j84u.cloudfront.net\/200,200\/image\/5566dc00f8b25e5b298b4568.png"}]}`,
		},
		"999": {
			statusCode:   http.StatusForbidden,
			responseData: `<?xml version="1.0" encoding="UTF-8"?><Error><Code>AccessDenied</Code><Message>Access Denied</Message><RequestId>9DF2F051BBC2C0C6</RequestId><HostId>tzjMqhELr2WzG7rlT01eXQhpR9Sx+UMVTkJ9Gc9J0ejESLv3fcokYp0GiEh5R7UFdzOxaeHTGMQ=</HostId></Error>`,
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlSlice := strings.Split(r.URL.Path, "/")
		id := urlSlice[len(urlSlice)-1]
		if response, ok := recipeResponseMap[id]; ok {
			w.WriteHeader(response.statusCode)
			_, _ = w.Write([]byte(response.responseData))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	testServer := httptest.NewServer(handler)
	ctx := context.Background()
	loggerCreator := func(_ context.Context) logrus.FieldLogger { return logger }
	httpClient := httpclient.New(10 * time.Second).WithRequestResponseLogger(loggerCreator)
	workerPool := worker.NewPool(10, loggerCreator)
	repo := repository.NewRecipeRepository(testServer.URL+"/", httpClient, workerPool)
	t.Run("all recipes exist", func(t *testing.T) {
		var buf bytes.Buffer
		logger.SetOutput(&buf)
		recipes, err := repo.FindByIDs(ctx, "1", "2", "3")
		if err != nil {
			t.Fatal("there shouldn't be error")
		}
		if len(recipes) != 3 {
			t.Fatal("there should be 3 recipes")
		}
		recipeIDs := make([]string, 0)
		for _, recipe := range recipes {
			recipeIDs = append(recipeIDs, recipe.Id)
		}
		sort.Strings(recipeIDs)
		if !reflect.DeepEqual(recipeIDs, []string{"1", "2", "3"}) {
			t.Fatal("recipe IDs should match")
		}
	})
	t.Run("one recipe doesn't exist", func(t *testing.T) {
		var buf bytes.Buffer
		logger.SetOutput(&buf)
		recipes, err := repo.FindByIDs(ctx, "1", "2", "999")
		if err != nil {
			t.Fatal("there shouldn't be error")
		}
		if !strings.Contains(buf.String(), "job return error: got 403 for URL") {
			t.Fatal("not found recipe should be reported in the log")
		}
		if len(recipes) != 2 {
			t.Fatal("there should be 2 recipes")
		}
		recipeIDs := make([]string, 0)
		for _, recipe := range recipes {
			recipeIDs = append(recipeIDs, recipe.Id)
		}
		sort.Strings(recipeIDs)
		if !reflect.DeepEqual(recipeIDs, []string{"1", "2"}) {
			t.Fatal("recipe IDs should match")
		}
	})
	t.Run("one recipe has error", func(t *testing.T) {
		var buf bytes.Buffer
		logger.SetOutput(&buf)
		recipes, err := repo.FindByIDs(ctx, "1", "2", "1000")
		if err == nil {
			t.Fatal("there should be error")
		}
		if !strings.Contains(err.Error(), "bad http status code") {
			t.Fatal("bad http status code should be report")
		}
		if len(recipes) != 0 {
			t.Fatal("there should be 0 recipes")
		}
	})
}
